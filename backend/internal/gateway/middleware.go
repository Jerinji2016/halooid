package gateway

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents the metrics for the API Gateway
type Metrics struct {
	RequestsTotal      *prometheus.CounterVec
	RequestDuration    *prometheus.HistogramVec
	RequestSize        *prometheus.SummaryVec
	ResponseSize       *prometheus.SummaryVec
	ResponseStatusCode *prometheus.CounterVec
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	m := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_requests_total",
				Help: "Total number of requests received by the API Gateway",
			},
			[]string{"method", "path", "service"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_request_duration_seconds",
				Help:    "Duration of requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "service"},
		),
		RequestSize: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "api_gateway_request_size_bytes",
				Help: "Size of requests in bytes",
			},
			[]string{"method", "path", "service"},
		),
		ResponseSize: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "api_gateway_response_size_bytes",
				Help: "Size of responses in bytes",
			},
			[]string{"method", "path", "service"},
		),
		ResponseStatusCode: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_response_status_code",
				Help: "Status code of responses",
			},
			[]string{"method", "path", "service", "status_code"},
		),
	}

	// Register metrics
	prometheus.MustRegister(
		m.RequestsTotal,
		m.RequestDuration,
		m.RequestSize,
		m.ResponseSize,
		m.ResponseStatusCode,
	)

	return m
}

// MetricsMiddleware creates a middleware that collects metrics
func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get service from path
		service := getServiceFromPath(r.URL.Path)

		// Start timer
		start := time.Now()

		// Create response writer wrapper
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call next handler
		next.ServeHTTP(ww, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		m.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, service).Inc()
		m.RequestDuration.WithLabelValues(r.Method, r.URL.Path, service).Observe(duration)
		m.RequestSize.WithLabelValues(r.Method, r.URL.Path, service).Observe(float64(r.ContentLength))
		m.ResponseSize.WithLabelValues(r.Method, r.URL.Path, service).Observe(float64(ww.BytesWritten()))
		m.ResponseStatusCode.WithLabelValues(r.Method, r.URL.Path, service, http.StatusText(ww.Status())).Inc()
	})
}

// getServiceFromPath extracts the service name from the path
func getServiceFromPath(path string) string {
	// Extract service name from path
	// Path format: /api/{service}/...
	if len(path) < 5 {
		return "unknown"
	}

	// Remove leading /api/
	path = path[5:]

	// Get service name
	for i, c := range path {
		if c == '/' {
			return path[:i]
		}
	}

	return path
}

// LoggingMiddleware creates a middleware that logs requests
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get request ID
			requestID := middleware.GetReqID(r.Context())

			// Log request
			logger.Printf("Request: %s %s %s", requestID, r.Method, r.URL.Path)

			// Create response writer wrapper
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Start timer
			start := time.Now()

			// Call next handler
			next.ServeHTTP(ww, r)

			// Log response
			duration := time.Since(start)
			logger.Printf("Response: %s %s %s %d %s", requestID, r.Method, r.URL.Path, ww.Status(), duration)
		})
	}
}
