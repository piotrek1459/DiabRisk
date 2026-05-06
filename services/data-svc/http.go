package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type dataServiceHandler struct {
	repository assessmentRepository
}

func newDataServiceHandler(repository assessmentRepository) dataServiceHandler {
	return dataServiceHandler{repository: repository}
}

func (h dataServiceHandler) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealthz)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/internal/assessments", h.handleCreateAssessment)
	mux.HandleFunc("/internal/users/", h.handleUserRoutes)
	return metricsMiddleware(mux)
}

func (h dataServiceHandler) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h dataServiceHandler) handleCreateAssessment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req createAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}

	record, err := h.repository.createAssessment(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, record)
}

func (h dataServiceHandler) handleUserRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	userID, ok := userIDFromAssessmentPath(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}

	records, err := h.repository.listAssessmentsByUser(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": records,
		"count": len(records),
	})
}

func userIDFromAssessmentPath(path string) (string, bool) {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 4 {
		return "", false
	}
	if parts[0] != "internal" || parts[1] != "users" || parts[3] != "assessments" {
		return "", false
	}
	if parts[2] == "" {
		return "", false
	}

	return parts[2], true
}

func runHTTPServer(port string, handler http.Handler) error {
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("data-svc HTTP API listening on :%s", port)
	return server.ListenAndServe()
}

func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
