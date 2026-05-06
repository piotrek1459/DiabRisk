package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "diabrisk_data_svc_http_requests_total",
			Help: "Total number of HTTP requests handled by data-svc.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "diabrisk_data_svc_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds for data-svc.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpResponseSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "diabrisk_data_svc_http_response_size_bytes",
			Help:    "HTTP response size in bytes for data-svc.",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000},
		},
		[]string{"method", "route", "status"},
	)
)

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	sizeBytes  int
}

func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *metricsResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	w.sizeBytes += size
	return size, err
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := newMetricsResponseWriter(w)

		next.ServeHTTP(recorder, r)

		route := dataServiceRouteName(r)
		status := strconv.Itoa(recorder.statusCode)
		labels := []string{r.Method, route, status}

		httpRequestsTotal.WithLabelValues(labels...).Inc()
		httpRequestDurationSeconds.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
		httpResponseSizeBytes.WithLabelValues(labels...).Observe(float64(recorder.sizeBytes))
	})
}

func dataServiceRouteName(r *http.Request) string {
	switch {
	case r.URL.Path == "/healthz":
		return "/healthz"
	case r.URL.Path == "/metrics":
		return "/metrics"
	case r.URL.Path == "/internal/assessments":
		return "/internal/assessments"
	case userIDFromAssessmentPathOnly(r.URL.Path):
		return "/internal/users/{user_id}/assessments"
	default:
		return "unmatched"
	}
}

func userIDFromAssessmentPathOnly(path string) bool {
	_, ok := userIDFromAssessmentPath(path)
	return ok
}
