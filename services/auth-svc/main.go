package main

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

var (
	pool *pgxpool.Pool
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     *string   `json:"full_name,omitempty"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://diabrisk:diabrisk_dev_password@postgres:5432/diabrisk?sslmode=disable")
	port := getEnv("PORT", "8081")
	adminEmail := getEnv("ADMIN_EMAIL", "admin@diabrisk.local")
	adminPassword := getEnv("ADMIN_PASSWORD", "default_admin_password")

	log.Println("Starting auth-svc...")
	log.Printf("Port: %s", port)

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// Seed default admin user with retry (waits for migrations to complete)
	go func() {
		maxRetries := 30
		for i := 0; i < maxRetries; i++ {
			if err := seedAdminUser(context.Background(), adminEmail, adminPassword); err != nil {
				log.Printf("Admin seed attempt %d/%d failed: %v", i+1, maxRetries, err)
				time.Sleep(2 * time.Second)
				continue
			}
			return
		}
		log.Printf("WARNING: Failed to seed admin user after %d attempts", maxRetries)
	}()

	r := gin.Default()
	r.Use(prometheusMiddleware())

	// CORS setup for local dev (Frontend on :5173, other services on k3d)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "http://diabrisk.local"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/auth/register", handleRegister)
	r.POST("/auth/login", handleLogin)
	r.POST("/auth/logout", handleLogout)
	r.GET("/auth/session", handleGetSession)

	log.Printf("✅ auth-svc listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user, err := createUser(context.Background(), req.Email, string(hashedPassword), req.FullName)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	session, sessionToken, err := createSession(context.Background(), user.ID)
	if err != nil {
		log.Printf("Failed to create session after registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.SetCookie(
		"session_token",
		sessionToken,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authenticateUser(context.Background(), req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to authenticate user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session, sessionToken, err := createSession(context.Background(), user.ID)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.SetCookie(
		"session_token",
		sessionToken,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "user": user})
}

func handleLogout(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Already logged out"})
		return
	}

	tokenHash := hashToken(sessionToken)
	_, err = pool.Exec(context.Background(),
		"UPDATE auth_sessions SET is_revoked = TRUE WHERE token_hash = $1",
		tokenHash,
	)
	if err != nil {
		log.Printf("Failed to revoke session: %v", err)
	}

	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func handleGetSession(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := validateSession(context.Background(), sessionToken)
	if err != nil {
		c.SetCookie("session_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func seedAdminUser(ctx context.Context, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (email, password_hash, full_name, role, last_login_at)
		VALUES ($1, $2, $3, 'admin', NOW())
		ON CONFLICT (email) DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			full_name = EXCLUDED.full_name,
			role = 'admin',
			last_login_at = NOW()
	`, email, string(hashedPassword), "Administrator")

	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	log.Printf("✅ Admin user seeded: %s", email)
	return nil
}

func createUser(ctx context.Context, email, passwordHash, fullName string) (*User, error) {
	var user User

	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, full_name, role, last_login_at)
		VALUES ($1, $2, $3, 'registered', NOW())
		RETURNING id, email, password_hash, full_name, role, created_at
	`, email, passwordHash, fullName).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func authenticateUser(ctx context.Context, email, password string) (*User, error) {
	var user User
	var passwordHash string

	err := pool.QueryRow(ctx, `
		SELECT id, email, password_hash, full_name, role, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName,
		&user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	_, _ = pool.Exec(ctx, "UPDATE users SET last_login_at = NOW() WHERE id = $1", user.ID)

	return &user, nil
}

func createSession(ctx context.Context, userID string) (*Session, string, error) {
	sessionToken := generateRandomString(64)
	tokenHash := hashToken(sessionToken)

	session := &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO auth_sessions (id, user_id, token_hash, expires_at, last_activity)
		VALUES ($1, $2, $3, $4, NOW())
	`, session.ID, session.UserID, session.TokenHash, session.ExpiresAt)

	if err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	return session, sessionToken, nil
}

func validateSession(ctx context.Context, sessionToken string) (*User, error) {
	tokenHash := hashToken(sessionToken)

	var user User
	err := pool.QueryRow(ctx, `
		SELECT u.id, u.email, u.password_hash, u.full_name, u.role, u.created_at
		FROM users u
		JOIN auth_sessions s ON u.id = s.user_id
		WHERE s.token_hash = $1
		  AND s.expires_at > NOW()
		  AND s.is_revoked = FALSE
	`, tokenHash).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid or expired session: %w", err)
	}

	_, _ = pool.Exec(ctx, "UPDATE auth_sessions SET last_activity = NOW() WHERE token_hash = $1", tokenHash)

	return &user, nil
}

func hashToken(token string) string {
	hash := sha512.Sum512([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
