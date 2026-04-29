package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "diabrisk_api_gateway_http_requests_total",
			Help: "Total number of HTTP requests handled by api-gateway.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "diabrisk_api_gateway_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds for api-gateway.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpResponseSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "diabrisk_api_gateway_http_response_size_bytes",
			Help:    "HTTP response size in bytes for api-gateway.",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000},
		},
		[]string{"method", "route", "status"},
	)
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := strconv.Itoa(c.Writer.Status())
		labels := []string{c.Request.Method, route, status}

		httpRequestsTotal.WithLabelValues(labels...).Inc()
		httpRequestDurationSeconds.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
		if size := c.Writer.Size(); size >= 0 {
			httpResponseSizeBytes.WithLabelValues(labels...).Observe(float64(size))
		}
	}
}
