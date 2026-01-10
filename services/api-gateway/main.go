package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS for local dev (Svelte on :5173)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// New /api/risk route to forward requests to ML service
	r.POST("/api/risk", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Marshal the request to JSON
		reqBody, err := json.Marshal(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		// Forward the request to the ML service
		resp, err := http.Post("http://65.109.169.137:8000/predict", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to ML service"})
			return
		}
		defer resp.Body.Close()

		// Read the response from the ML service
		var mlResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&mlResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ML service response"})
			return
		}

		c.JSON(http.StatusOK, mlResponse)
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
