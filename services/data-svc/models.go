package main

import "time"

type createAssessmentRequest struct {
	UserID      string                 `json:"user_id"`
	Features    map[string]interface{} `json:"features"`
	RiskPercent float64                `json:"risk_percent"`
	Category    string                 `json:"category"`
	GeneratedAt time.Time              `json:"generated_at"`
}

type assessmentRecord struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id"`
	ModelVersionID string                 `json:"model_version_id"`
	Features       map[string]interface{} `json:"features"`
	RiskPercent    float64                `json:"risk_percent"`
	Category       string                 `json:"category"`
	Message        string                 `json:"message"`
	CreatedAt      time.Time              `json:"created_at"`
}
