package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func withAuthServiceURL(t *testing.T, url string) {
	t.Helper()
	t.Setenv("AUTH_SERVICE_URL", url)
}

func decodeJSONBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode JSON response %q: %v", w.Body.String(), err)
	}
	return body
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

// --- auth middleware tests ---

func TestAuthMiddleware_NoCookieReturnsUnauthorized(t *testing.T) {
	r := gin.New()
	r.Use(authMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "allowed"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"Not authenticated"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_RejectsInvalidSession(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET to auth service, got %s", r.Method)
		}
		if r.URL.Path != "/auth/session" {
			t.Errorf("expected /auth/session, got %s", r.URL.Path)
		}
		if cookie, err := r.Cookie("session_token"); err != nil || cookie.Value != "bad-token" {
			t.Errorf("expected session_token cookie to be forwarded, got cookie=%v err=%v", cookie, err)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer authServer.Close()
	withAuthServiceURL(t, authServer.URL)

	r := gin.New()
	r.Use(authMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "allowed"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "bad-token"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"Invalid session"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_AllowsValidSessionAndStoresUser(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET to auth service, got %s", r.Method)
		}
		if r.URL.Path != "/auth/session" {
			t.Errorf("expected /auth/session, got %s", r.URL.Path)
		}
		if cookie, err := r.Cookie("session_token"); err != nil || cookie.Value != "good-token" {
			t.Errorf("expected session_token cookie to be forwarded, got cookie=%v err=%v", cookie, err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"user-1","email":"test@example.com"}`)
	}))
	defer authServer.Close()
	withAuthServiceURL(t, authServer.URL)

	r := gin.New()
	r.Use(authMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "missing user"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "good-token"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", w.Code, w.Body.String())
	}
	body := decodeJSONBody(t, w)
	if body["id"] != "user-1" || body["email"] != "test@example.com" {
		t.Errorf("unexpected user body: %#v", body)
	}
}

// --- proxy tests ---

func TestProxyToAuthService_ForwardsRequestAndCopiesResponse(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/auth/login" {
			t.Errorf("expected /auth/login, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("next"); got != "/dashboard" {
			t.Errorf("expected next=/dashboard query, got %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected content type application/json, got %q", got)
		}
		if cookie, err := r.Cookie("session_token"); err != nil || cookie.Value != "existing-token" {
			t.Errorf("expected session_token cookie to be forwarded, got cookie=%v err=%v", cookie, err)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read proxied body: %v", err)
			return
		}
		if string(body) != `{"email":"user@example.com"}` {
			t.Errorf("unexpected proxied body: %s", string(body))
		}

		w.Header().Set("X-Auth-Result", "forwarded")
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "new-token", Path: "/", HttpOnly: true})
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer authServer.Close()
	withAuthServiceURL(t, authServer.URL)

	r := gin.New()
	r.POST("/auth/login", proxyToAuthService)

	req := httptest.NewRequest(http.MethodPost, "/auth/login?next=%2Fdashboard", strings.NewReader(`{"email":"user@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "existing-token"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	if w.Body.String() != `{"ok":true}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
	if got := w.Header().Get("X-Auth-Result"); got != "forwarded" {
		t.Errorf("expected X-Auth-Result header to be copied, got %q", got)
	}
	if got := w.Header().Get("Set-Cookie"); !strings.Contains(got, "session_token=new-token") {
		t.Errorf("expected Set-Cookie to be copied, got %q", got)
	}
}

func TestProxyToAuthService_Returns500WhenAuthServiceUnavailable(t *testing.T) {
	withAuthServiceURL(t, "http://127.0.0.1:1")

	r := gin.New()
	r.POST("/auth/login", proxyToAuthService)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com"}`))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"Failed to reach auth service"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}
