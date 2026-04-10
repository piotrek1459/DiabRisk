package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- getEnv tests ---

func TestGetEnv_WithDefault(t *testing.T) {
	val := getEnv("DIABRISK_TEST_NONEXISTENT_KEY_99999", "default-val")
	if val != "default-val" {
		t.Errorf("expected 'default-val', got '%s'", val)
	}
}

func TestGetEnv_WithValue(t *testing.T) {
	t.Setenv("DIABRISK_TEST_KEY", "real-value")
	val := getEnv("DIABRISK_TEST_KEY", "default-val")
	if val != "real-value" {
		t.Errorf("expected 'real-value', got '%s'", val)
	}
}

func TestGetEnv_EmptyStringUsesDefault(t *testing.T) {
	t.Setenv("DIABRISK_TEST_EMPTY", "")
	val := getEnv("DIABRISK_TEST_EMPTY", "fallback")
	if val != "fallback" {
		t.Errorf("expected 'fallback' for empty env var, got '%s'", val)
	}
}

// --- healthz endpoint test ---

func setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}

func TestHealthzEndpoint(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", body)
	}
}
