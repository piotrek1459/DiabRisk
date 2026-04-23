package main

import "time"

type riskPredictionRequest struct {
	Features map[string]interface{} `json:"features"`
}

type riskPredictionResult struct {
	AssessmentID string  `json:"AssessmentID,omitempty"`
	RiskPercent  float64 `json:"RiskPercent"`
	Category     string  `json:"Category"`
	Message      string  `json:"Message"`
	GeneratedAt  string  `json:"GeneratedAt,omitempty"`
}

type assessmentCandidate struct {
	UserID      string                 `json:"user_id"`
	Features    map[string]interface{} `json:"features"`
	RiskPercent float64                `json:"risk_percent"`
	Category    string                 `json:"category"`
	GeneratedAt time.Time              `json:"generated_at"`
}

type assessmentHistoryItem struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id,omitempty"`
	ModelVersionID string                 `json:"model_version_id,omitempty"`
	Features       map[string]interface{} `json:"features"`
	RiskPercent    float64                `json:"risk_percent"`
	Category       string                 `json:"category"`
	Message        string                 `json:"message"`
	CreatedAt      time.Time              `json:"created_at"`
}

type assessmentHistoryResponse struct {
	Items []assessmentHistoryItem `json:"items"`
	Count int                     `json:"count"`
}

type authUser struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FullName  *string    `json:"full_name,omitempty"`
	Role      string     `json:"role"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
