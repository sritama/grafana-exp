package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"runtime"
	"time"
)

// Custom metrics
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	goroutinesGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Current number of goroutines",
		},
	)

	memoryAllocGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Current memory allocation in bytes",
		},
	)

	memoryHeapGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_heap_bytes",
			Help: "Current heap memory usage in bytes",
		},
	)

	gcDurationGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "gc_duration_seconds",
			Help: "Duration of the last garbage collection cycle",
		},
	)
)

// Middleware to collect HTTP metrics
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()

		// Record metrics
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", wrapped.statusCode)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

// Response writer wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Function to collect runtime metrics
func collectRuntimeMetrics() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Goroutines
		goroutinesGauge.Set(float64(runtime.NumGoroutine()))

		// Memory stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		memoryAllocGauge.Set(float64(m.Alloc))
		memoryHeapGauge.Set(float64(m.HeapAlloc))
		gcDurationGauge.Set(float64(m.PauseNs[(m.NumGC+255)%256]) / 1e9) // Convert nanoseconds to seconds
	}
}
