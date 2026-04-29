package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// authMiddlewareWithConfig validates the session by calling auth-svc.
func authMiddlewareWithConfig(config serviceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// Validate session with auth-svc
		req, _ := http.NewRequest("GET", config.authServiceURL+"/auth/session", nil)
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
		var user authUser
		if err := json.NewDecoder(resp.Body).Decode(&user); err == nil {
			c.Set("user", user)
		}

		c.Next()
	}
}

func main() {
	config := newServiceConfig()
	r := createRouterWithConfig(config)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func createRouter() *gin.Engine {
	return createRouterWithConfig(newServiceConfig())
}

func createRouterWithConfig(config serviceConfig) *gin.Engine {
	r := gin.Default()
	handler := newGatewayHandler(config)
	r.Use(prometheusMiddleware())

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
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Auth routes - proxy to auth-svc (no auth required)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", proxyToAuthServiceWithConfig(config))
		authRoutes.POST("/login", proxyToAuthServiceWithConfig(config))
		authRoutes.POST("/logout", proxyToAuthServiceWithConfig(config))
		authRoutes.GET("/session", proxyToAuthServiceWithConfig(config))
	}

	// Protected API routes
	api := r.Group("/api")
	api.Use(authMiddlewareWithConfig(config))
	{
		api.POST("/risk", handler.handleRiskPrediction)
		api.GET("/features", handler.handleFeatures)
		api.GET("/history", handler.handleAssessmentHistory)
	}

	return r
}

// proxyToAuthService forwards requests to auth-svc, preserving cookies.
func proxyToAuthService(c *gin.Context) {
	proxyToAuthServiceWithConfig(newServiceConfig())(c)
}

func proxyToAuthServiceWithConfig(config serviceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL := config.authServiceURL + c.Request.URL.Path
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
}
