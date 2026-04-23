package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func withMLServiceURL(t *testing.T, url string) {
	t.Helper()
	original := mlServiceURL
	mlServiceURL = url
	t.Cleanup(func() {
		mlServiceURL = original
	})
}

func withDataServiceURL(t *testing.T, url string) {
	t.Helper()
	original := dataServiceURL
	dataServiceURL = url
	t.Cleanup(func() {
		dataServiceURL = original
	})
}

func TestGatewayIntegration_ProtectedRiskRouteValidatesSessionAndCallsML(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/auth/session" {
			t.Fatalf("unexpected auth request: %s %s", r.Method, r.URL.Path)
		}
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value != "valid-session" {
			t.Fatalf("expected forwarded session cookie, got cookie=%v err=%v", cookie, err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"user-1","email":"test@example.com"}`)
	}))
	defer authServer.Close()

	mlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/predict" {
			t.Fatalf("unexpected ml request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode ml request: %v", err)
		}
		if body["features"]["BMI"] != 30.0 {
			t.Fatalf("expected BMI feature to be forwarded, got %#v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"RiskPercent":0.63,"Category":"medium","Message":"ok"}`)
	}))
	defer mlServer.Close()

	dataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/internal/assessments" {
			t.Fatalf("unexpected data service request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"assessment-1","features":{"BMI":30},"risk_percent":0.63,"category":"medium","message":"ok","created_at":"2026-04-23T12:00:00Z"}`)
	}))
	defer dataServer.Close()

	withAuthServiceURL(t, authServer.URL)
	withMLServiceURL(t, mlServer.URL)
	withDataServiceURL(t, dataServer.URL)

	router := createRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/risk", strings.NewReader(`{"features":{"BMI":30}}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "valid-session"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", w.Code, w.Body.String())
	}

	body := decodeJSONBody(t, w)
	if body["Category"] != "medium" || body["RiskPercent"] != 0.63 {
		t.Fatalf("unexpected gateway risk response: %#v", body)
	}
}

func TestGatewayIntegration_FeaturesRouteRequiresAuthAndProxiesMLResponse(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"user-1","email":"test@example.com"}`)
	}))
	defer authServer.Close()

	mlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/features" {
			t.Fatalf("unexpected ml request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"feature_names":["BMI","Age"],"count":2}`)
	}))
	defer mlServer.Close()

	withAuthServiceURL(t, authServer.URL)
	withMLServiceURL(t, mlServer.URL)

	router := createRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/features", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "valid-session"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", w.Code, w.Body.String())
	}

	body := decodeJSONBody(t, w)
	if body["count"] != float64(2) {
		t.Fatalf("unexpected features response: %#v", body)
	}
}
