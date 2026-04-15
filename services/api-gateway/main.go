package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	authServiceURL = getEnv("AUTH_SERVICE_URL", "http://auth-svc:8081")
	mlServiceURL   = getEnv("ML_SERVICE_URL", "http://ml-api:8000")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// authMiddleware validates the session by calling auth-svc
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// Validate session with auth-svc
		req, _ := http.NewRequest("GET", authServiceURL+"/auth/session", nil)
		req.AddCookie(&http.Cookie{Name: "session_token", Value: cookie})

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Parse user info from response
		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err == nil {
			c.Set("user", userInfo)
		}

		c.Next()
	}
}

func main() {
	r := createRouter()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func createRouter() *gin.Engine {
	r := gin.Default()

	// CORS for local dev (Svelte on :5173)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://diabrisk.local"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Auth routes - proxy to auth-svc (no auth required)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", proxyToAuthService)
		authRoutes.POST("/login", proxyToAuthService)
		authRoutes.POST("/logout", proxyToAuthService)
		authRoutes.GET("/session", proxyToAuthService)
	}

	// Protected API routes
	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		// Risk assessment route (protected)
		api.POST("/risk", func(c *gin.Context) {
			var req map[string]interface{}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			featuresPayload := req
			if nested, ok := req["features"].(map[string]interface{}); ok {
				featuresPayload = nested
			}

			// Build ML API request
			mlReq := map[string]interface{}{
				"features": featuresPayload,
			}

			// Marshal the request to JSON
			reqBody, err := json.Marshal(mlReq)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
				return
			}

			// Forward the request to the ML service
			mlResp, err := http.Post(mlServiceURL+"/predict", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to ML service"})
				return
			}
			defer mlResp.Body.Close()

			// Read the response from the ML service
			var mlResponse map[string]interface{}
			if err := json.NewDecoder(mlResp.Body).Decode(&mlResponse); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ML service response"})
				return
			}

			c.JSON(http.StatusOK, mlResponse)
		})

		// Features endpoint (list available features)
		api.GET("/features", func(c *gin.Context) {
			mlResp, err := http.Get(mlServiceURL + "/features")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to ML service"})
				return
			}
			defer mlResp.Body.Close()

			var features map[string]interface{}
			if err := json.NewDecoder(mlResp.Body).Decode(&features); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ML service response"})
				return
			}

			c.JSON(http.StatusOK, features)
		})
	}

	return r
}

// proxyToAuthService forwards requests to auth-svc, preserving cookies
func proxyToAuthService(c *gin.Context) {
	targetURL := authServiceURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers (this includes Cookie header)
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Forward request
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach auth service"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers (including Set-Cookie)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Copy status code
	c.Status(resp.StatusCode)

	// Copy body
	io.Copy(c.Writer, resp.Body)
}
