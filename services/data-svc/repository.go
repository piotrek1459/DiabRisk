package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type assessmentRepository struct {
	pool *pgxpool.Pool
}

const defaultModelVersion = "v1.0.0"

func newAssessmentRepository(pool *pgxpool.Pool) assessmentRepository {
	return assessmentRepository{pool: pool}
}

func (r assessmentRepository) createAssessment(ctx context.Context, req createAssessmentRequest) (*assessmentRecord, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if len(req.Features) == 0 {
		return nil, fmt.Errorf("features are required")
	}

	modelVersionID, err := r.defaultModelVersionID(ctx)
	if err != nil {
		return nil, err
	}

	featuresJSON, err := json.Marshal(req.Features)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal features: %w", err)
	}

	createdAt := req.GeneratedAt.UTC()
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	category := normalizeCategory(req.Category)
	riskLevel, err := categoryToRiskLevel(category)
	if err != nil {
		return nil, err
	}

	record := &assessmentRecord{
		UserID:         req.UserID,
		ModelVersionID: modelVersionID,
		Features:       req.Features,
		RiskPercent:    req.RiskPercent,
		Category:       category,
		Message:        messageForCategory(category),
		CreatedAt:      createdAt,
	}

	err = r.pool.QueryRow(ctx, `
		INSERT INTO assessments (
			user_id,
			model_version_id,
			features,
			raw_score,
			risk_level,
			created_at
		)
		VALUES ($1, $2, $3::jsonb, $4, $5, $6)
		RETURNING id
	`,
		record.UserID,
		record.ModelVersionID,
		string(featuresJSON),
		record.RiskPercent,
		riskLevel,
		record.CreatedAt,
	).Scan(&record.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert assessment: %w", err)
	}

	return record, nil
}

func (r assessmentRepository) listAssessmentsByUser(ctx context.Context, userID string) ([]assessmentRecord, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, model_version_id, features, raw_score, risk_level, created_at
		FROM assessments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assessments: %w", err)
	}
	defer rows.Close()

	records := make([]assessmentRecord, 0)
	for rows.Next() {
		var (
			record       assessmentRecord
			featuresJSON []byte
			riskLevel    string
		)

		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.ModelVersionID,
			&featuresJSON,
			&record.RiskPercent,
			&riskLevel,
			&record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan assessment: %w", err)
		}

		if err := json.Unmarshal(featuresJSON, &record.Features); err != nil {
			return nil, fmt.Errorf("failed to unmarshal assessment features: %w", err)
		}

		record.Category = riskLevelToCategory(riskLevel)
		record.Message = messageForCategory(record.Category)
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("assessment rows error: %w", err)
	}

	return records, nil
}

func (r assessmentRepository) defaultModelVersionID(ctx context.Context) (string, error) {
	var modelVersionID string
	err := r.pool.QueryRow(ctx, `
		SELECT id
		FROM model_versions
		WHERE version = $1
	`, defaultModelVersion).Scan(&modelVersionID)
	if err != nil {
		return "", fmt.Errorf("failed to resolve active model version: %w", err)
	}

	return modelVersionID, nil
}

func normalizeCategory(category string) string {
	return strings.ToLower(strings.TrimSpace(category))
}

func categoryToRiskLevel(category string) (string, error) {
	switch normalizeCategory(category) {
	case "low":
		return "no_diabetes", nil
	case "medium":
		return "prediabetes", nil
	case "high":
		return "diabetes", nil
	default:
		return "", fmt.Errorf("unsupported category %q", category)
	}
}

func riskLevelToCategory(riskLevel string) string {
	switch strings.ToLower(strings.TrimSpace(riskLevel)) {
	case "no_diabetes":
		return "low"
	case "prediabetes":
		return "medium"
	case "diabetes":
		return "high"
	default:
		return "unknown"
	}
}

func messageForCategory(category string) string {
	switch normalizeCategory(category) {
	case "high":
		return "High risk detected. Immediate medical consultation recommended."
	case "medium":
		return "Moderate risk detected. Consider scheduling a medical checkup."
	case "low":
		return "Low risk detected. No immediate action required."
	default:
		return "Risk assessment completed."
	}
}
