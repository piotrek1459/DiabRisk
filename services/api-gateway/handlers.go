package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type assessmentRecorder interface {
	PrepareForHistory(*assessmentCandidate) error
}

type noopAssessmentRecorder struct{}

func (noopAssessmentRecorder) PrepareForHistory(*assessmentCandidate) error {
	return nil
}

type gatewayHandler struct {
	config   serviceConfig
	recorder assessmentRecorder
}

func newGatewayHandler(config serviceConfig) gatewayHandler {
	return gatewayHandler{
		config:   config,
		recorder: noopAssessmentRecorder{},
	}
}

func authMiddleware() gin.HandlerFunc {
	return authMiddlewareWithConfig(loadServiceConfig())
}

func (h gatewayHandler) handleRiskPrediction(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	featuresPayload := extractFeaturesPayload(req)
	mlReq := riskPredictionRequest{Features: featuresPayload}
	reqBody, err := json.Marshal(mlReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	mlResp, err := http.Post(h.config.mlServiceURL+"/predict", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to ML service"})
		return
	}
	defer mlResp.Body.Close()

	var prediction riskPredictionResult
	if err := json.NewDecoder(mlResp.Body).Decode(&prediction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ML service response"})
		return
	}

	prediction.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	user, err := currentAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authenticated user context"})
		return
	}

	assessment := &assessmentCandidate{
		UserID:      user.ID,
		Features:    featuresPayload,
		RiskPercent: prediction.RiskPercent,
		Category:    prediction.Category,
		Message:     prediction.Message,
		GeneratedAt: time.Now().UTC(),
	}

	if err := h.recorder.PrepareForHistory(assessment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare assessment for history"})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

func (h gatewayHandler) handleFeatures(c *gin.Context) {
	mlResp, err := http.Get(h.config.mlServiceURL + "/features")
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
}

func extractFeaturesPayload(req map[string]interface{}) map[string]interface{} {
	if nested, ok := req["features"].(map[string]interface{}); ok {
		return nested
	}

	return req
}

func currentAuthUser(c *gin.Context) (*authUser, error) {
	rawUser, exists := c.Get("user")
	if !exists {
		return nil, errors.New("missing user in request context")
	}

	userMap, ok := rawUser.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected user payload type")
	}

	userID, ok := userMap["id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("missing user id")
	}

	email, _ := userMap["email"].(string)
	role, _ := userMap["role"].(string)
	fullName, _ := userMap["full_name"].(string)

	var fullNamePtr *string
	if fullName != "" {
		fullNamePtr = &fullName
	}

	return &authUser{
		ID:       userID,
		Email:    email,
		FullName: fullNamePtr,
		Role:     role,
	}, nil
}
