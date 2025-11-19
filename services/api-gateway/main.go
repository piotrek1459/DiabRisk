package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RiskRequest struct {
	Age        int     `json:"age" binding:"required"`
	BMI        float64 `json:"bmi" binding:"required"`
	SystolicBP int     `json:"systolic_bp" binding:"required"`
	Glucose    int     `json:"glucose" binding:"required"`
	Smoker     bool    `json:"smoker"`
}

type RiskResponse struct {
	RiskPercent float64 `json:"risk_percent"`
	Category    string  `json:"category"` // "low", "medium", "high"
	Message     string  `json:"message"`
}

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

	// For now: dummy risk logic
	r.POST("/api/risk", func(c *gin.Context) {
		var req RiskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// very rough & fake scoring – just so frontend has something
		score := 0.0
		if req.Age > 45 {
			score += 20
		}
		if req.BMI >= 30 {
			score += 30
		}
		if req.SystolicBP >= 140 {
			score += 20
		}
		if req.Glucose >= 126 {
			score += 25
		}
		if req.Smoker {
			score += 10
		}
		if score > 100 {
			score = 100
		}

		category := "low"
		message := "Estimated low risk. Keep up healthy habits."
		if score >= 33 && score < 66 {
			category = "medium"
			message = "Moderate risk. Consider lifestyle improvements and consulting a doctor."
		} else if score >= 66 {
			category = "high"
			message = "High estimated risk. Please consult a medical professional."
		}

		c.JSON(http.StatusOK, RiskResponse{
			RiskPercent: score,
			Category:    category,
			Message:     message,
		})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
