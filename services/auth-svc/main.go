package main

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	pool        *pgxpool.Pool
	oauthConfig *oauth2.Config
)

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	GoogleID   *string   `json:"google_id,omitempty"`
	FullName   *string   `json:"full_name,omitempty"`
	PictureURL *string   `json:"picture_url,omitempty"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type Session struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://diabrisk:diabrisk_dev_password@postgres:5432/diabrisk?sslmode=disable")
	port := getEnv("PORT", "8081")
	clientID := getEnv("GOOGLE_CLIENT_ID", "")
	clientSecret := getEnv("GOOGLE_CLIENT_SECRET", "")
	redirectURL := getEnv("REDIRECT_URL", "http://diabrisk.local/auth/google/callback")

	log.Println("Starting auth-svc...")
	log.Printf("Port: %s", port)
	log.Printf("Redirect URL: %s", redirectURL)

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

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

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/auth/google/login", handleGoogleLogin)
	r.GET("/auth/google/callback", handleGoogleCallback)
	r.POST("/auth/logout", handleLogout)
	r.GET("/auth/session", handleGetSession)

	log.Printf("✅ auth-svc listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleGoogleLogin(c *gin.Context) {
	state := generateRandomString(32)
	c.SetCookie("oauth_state", state, 600, "/", "localhost", false, true)
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || state != cookieState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "localhost", false, true)

	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	user, err := findOrCreateUser(context.Background(), googleUser)
	if err != nil {
		log.Printf("Failed to find/create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
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
		"localhost",
		false,
		true,
	)

	c.Redirect(http.StatusTemporaryRedirect, "/")
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

func findOrCreateUser(ctx context.Context, googleUser GoogleUserInfo) (*User, error) {
	var user User

	err := pool.QueryRow(ctx, `
		SELECT id, email, google_id, full_name, picture_url, role, created_at
		FROM users
		WHERE google_id = $1
	`, googleUser.ID).Scan(
		&user.ID, &user.Email, &user.GoogleID, &user.FullName,
		&user.PictureURL, &user.Role, &user.CreatedAt,
	)

	if err == nil {
		_, _ = pool.Exec(ctx, "UPDATE users SET last_login_at = NOW() WHERE id = $1", user.ID)
		return &user, nil
	}

	log.Printf("Creating new user: %s (%s)", googleUser.Email, googleUser.Name)

	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, google_id, full_name, picture_url, role, last_login_at)
		VALUES ($1, $2, $3, $4, 'registered', NOW())
		RETURNING id, email, google_id, full_name, picture_url, role, created_at
	`, googleUser.Email, googleUser.ID, googleUser.Name, googleUser.Picture).Scan(
		&user.ID, &user.Email, &user.GoogleID, &user.FullName,
		&user.PictureURL, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

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
		SELECT u.id, u.email, u.google_id, u.full_name, u.picture_url, u.role, u.created_at
		FROM users u
		JOIN auth_sessions s ON u.id = s.user_id
		WHERE s.token_hash = $1
		  AND s.expires_at > NOW()
		  AND s.is_revoked = FALSE
	`, tokenHash).Scan(
		&user.ID, &user.Email, &user.GoogleID, &user.FullName,
		&user.PictureURL, &user.Role, &user.CreatedAt,
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
