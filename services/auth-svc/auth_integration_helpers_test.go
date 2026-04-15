//go:build integration

package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ginRouterForAuthIntegration() *gin.Engine {
	r := gin.New()
	r.POST("/auth/register", handleRegister)
	r.POST("/auth/login", handleLogin)
	r.POST("/auth/logout", handleLogout)
	r.GET("/auth/session", handleGetSession)
	return r
}

func firstCookie(t testingT, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("missing cookie %q in %v", name, w.Result().Cookies())
	return nil
}

func httptestRequestWithCookie(method, path string, cookie *http.Cookie) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.AddCookie(cookie)
	return req
}

func serveRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type testingT interface {
	Helper()
	Fatalf(format string, args ...interface{})
}
