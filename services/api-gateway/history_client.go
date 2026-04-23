package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type assessmentRecorder interface {
	Record(context.Context, *assessmentCandidate) (*assessmentHistoryItem, error)
	ListByUser(context.Context, string) ([]assessmentHistoryItem, error)
}

type dataServiceRecorder struct {
	baseURL string
	client  *http.Client
}

func newAssessmentRecorder(config serviceConfig) assessmentRecorder {
	return dataServiceRecorder{
		baseURL: strings.TrimRight(config.dataServiceURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (r dataServiceRecorder) Record(ctx context.Context, candidate *assessmentCandidate) (*assessmentHistoryItem, error) {
	payload, err := json.Marshal(candidate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode assessment candidate: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/internal/assessments", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create data service request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create assessment in data service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, recorderResponseError("data service rejected assessment write", resp)
	}

	var item assessmentHistoryItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode assessment response: %w", err)
	}

	return &item, nil
}

func (r dataServiceRecorder) ListByUser(ctx context.Context, userID string) ([]assessmentHistoryItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/internal/users/"+userID+"/assessments", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create history request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assessment history from data service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, recorderResponseError("data service rejected assessment history read", resp)
	}

	var history assessmentHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("failed to decode assessment history response: %w", err)
	}

	return history.Items, nil
}

func recorderResponseError(prefix string, resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		return fmt.Errorf("%s: status %d", prefix, resp.StatusCode)
	}

	return fmt.Errorf("%s: status %d: %s", prefix, resp.StatusCode, strings.TrimSpace(string(body)))
}
