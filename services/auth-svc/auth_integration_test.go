//go:build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func requireAuthIntegrationDB(t *testing.T) {
	t.Helper()
	dbURL := getEnv("AUTH_INTEGRATION_DATABASE_URL", "")
	if dbURL == "" {
		dbURL = getEnv("DATABASE_URL", "")
	}
	if dbURL == "" {
		t.Skip("set AUTH_INTEGRATION_DATABASE_URL or DATABASE_URL to run auth integration tests")
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
		pool = nil
	})

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}

func TestAuthIntegration_RegisterLoginSessionLogout(t *testing.T) {
	requireAuthIntegrationDB(t)

	router := ginRouterForAuthIntegration()
	email := fmt.Sprintf("go-auth-%d@example.test", time.Now().UnixNano())
	password := "integration-password-123"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE email = $1", email)
	})

	register := performJSONRequest(
		router,
		http.MethodPost,
		"/auth/register",
		fmt.Sprintf(`{"email":%q,"password":%q,"full_name":"Go Integration"}`, email, password),
	)
	if register.Code != http.StatusCreated {
		t.Fatalf("expected register 201, got %d with body %s", register.Code, register.Body.String())
	}
	if !strings.Contains(register.Header().Get("Set-Cookie"), "session_token=") {
		t.Fatalf("expected register to set session cookie, got %q", register.Header().Get("Set-Cookie"))
	}

	login := performJSONRequest(
		router,
		http.MethodPost,
		"/auth/login",
		fmt.Sprintf(`{"email":%q,"password":%q}`, email, password),
	)
	if login.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d with body %s", login.Code, login.Body.String())
	}

	sessionCookie := firstCookie(t, login, "session_token")
	sessionReq := httptestRequestWithCookie(http.MethodGet, "/auth/session", sessionCookie)
	session := serveRequest(router, sessionReq)
	if session.Code != http.StatusOK {
		t.Fatalf("expected session 200, got %d with body %s", session.Code, session.Body.String())
	}
	if !strings.Contains(session.Body.String(), email) {
		t.Fatalf("expected session response for %s, got %s", email, session.Body.String())
	}

	logoutReq := httptestRequestWithCookie(http.MethodPost, "/auth/logout", sessionCookie)
	logout := serveRequest(router, logoutReq)
	if logout.Code != http.StatusOK {
		t.Fatalf("expected logout 200, got %d with body %s", logout.Code, logout.Body.String())
	}

	afterLogoutReq := httptestRequestWithCookie(http.MethodGet, "/auth/session", sessionCookie)
	afterLogout := serveRequest(router, afterLogoutReq)
	if afterLogout.Code != http.StatusUnauthorized {
		t.Fatalf("expected session 401 after logout, got %d with body %s", afterLogout.Code, afterLogout.Body.String())
	}
}
