package main

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- hashToken tests ---

func TestHashToken_Deterministic(t *testing.T) {
	token := "test-session-token-abc123"
	hash1 := hashToken(token)
	hash2 := hashToken(token)
	if hash1 != hash2 {
		t.Errorf("hashToken is not deterministic: %s != %s", hash1, hash2)
	}
}

func TestHashToken_ProducesSHA512(t *testing.T) {
	hash := hashToken("hello")
	// SHA-512 produces 128 hex characters
	if len(hash) != 128 {
		t.Errorf("expected 128 hex chars, got %d", len(hash))
	}
	// Must be valid hex
	if _, err := hex.DecodeString(hash); err != nil {
		t.Errorf("hash is not valid hex: %v", err)
	}
}

func TestHashToken_DifferentInputsDifferentHashes(t *testing.T) {
	h1 := hashToken("token-a")
	h2 := hashToken("token-b")
	if h1 == h2 {
		t.Error("different inputs produced the same hash")
	}
}

// --- generateRandomString tests ---

func TestGenerateRandomString_Length(t *testing.T) {
	for _, length := range []int{16, 32, 64} {
		s := generateRandomString(length)
		if len(s) != length {
			t.Errorf("generateRandomString(%d) returned string of length %d", length, len(s))
		}
	}
}

func TestGenerateRandomString_Unique(t *testing.T) {
	s1 := generateRandomString(64)
	s2 := generateRandomString(64)
	if s1 == s2 {
		t.Error("two calls to generateRandomString returned the same value")
	}
}

// --- getEnv tests ---

func TestGetEnv_WithDefault(t *testing.T) {
	val := getEnv("DIABRISK_TEST_NONEXISTENT_KEY_12345", "fallback")
	if val != "fallback" {
		t.Errorf("expected 'fallback', got '%s'", val)
	}
}

func TestGetEnv_WithValue(t *testing.T) {
	t.Setenv("DIABRISK_TEST_KEY", "myvalue")
	val := getEnv("DIABRISK_TEST_KEY", "fallback")
	if val != "myvalue" {
		t.Errorf("expected 'myvalue', got '%s'", val)
	}
}

// --- HTTP handler tests (no DB required) ---

func TestHandleRegister_InvalidJSON(t *testing.T) {
	r := gin.New()
	r.POST("/auth/register", handleRegister)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleLogin_InvalidJSON(t *testing.T) {
	r := gin.New()
	r.POST("/auth/login", handleLogin)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleRegister_MissingFields(t *testing.T) {
	r := gin.New()
	r.POST("/auth/register", handleRegister)

	// Missing password
	body := `{"email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing password, got %d", w.Code)
	}
}

func TestHandleRegister_PasswordTooShort(t *testing.T) {
	r := gin.New()
	r.POST("/auth/register", handleRegister)

	body := `{"email": "test@example.com", "password": "abc"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for short password, got %d", w.Code)
	}
}

func TestHandleLogin_MissingEmail(t *testing.T) {
	r := gin.New()
	r.POST("/auth/login", handleLogin)

	body := `{"password": "test123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing email, got %d", w.Code)
	}
}
