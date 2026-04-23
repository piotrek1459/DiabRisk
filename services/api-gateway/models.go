package main

import "time"

type riskPredictionRequest struct {
	Features map[string]interface{} `json:"features"`
}

type riskPredictionResult struct {
	RiskPercent float64 `json:"RiskPercent"`
	Category    string  `json:"Category"`
	Message     string  `json:"Message"`
	GeneratedAt string  `json:"GeneratedAt,omitempty"`
}

type assessmentCandidate struct {
	UserID      string                 `json:"user_id"`
	Features    map[string]interface{} `json:"features"`
	RiskPercent float64                `json:"risk_percent"`
	Category    string                 `json:"category"`
	Message     string                 `json:"message"`
	GeneratedAt time.Time              `json:"generated_at"`
}

type authUser struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FullName  *string    `json:"full_name,omitempty"`
	Role      string     `json:"role"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
