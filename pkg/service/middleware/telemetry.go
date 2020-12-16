package middleware

import (
	"fmt"
	"net/http"
	"time"

	mw "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds",
			Help:       "HTTP request duration summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.05, 0.95: 0.05, 0.99: 0.05},
		},
		[]string{"handler", "verb", "status"},
	)

	RequestErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "total number of valid HTTP requests that could not be served",
		},
		[]string{"handler", "verb", "status"},
	)

	RequestSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_size_bytes",
			Help:       "The request size of HTTP requests in bytes",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.05, 0.95: 0.05, 0.99: 0.05},
		},
	)

	ResponseSize = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_response_size_bytes",
			Help:       "The response size of HTTP requests in bytes",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.05, 0.95: 0.05, 0.99: 0.05},
		},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestErrors)
	prometheus.MustRegister(RequestSize)
	prometheus.MustRegister(ResponseSize)
}

func PrometheusTelemetry(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := mw.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(start)
		status := fmt.Sprintf("%d", ww.Status())
		RequestDuration.WithLabelValues(r.URL.Path, r.Method, status).Observe((duration.Seconds()))
		RequestSize.Observe(float64(r.ContentLength))
		ResponseSize.Observe(float64(ww.BytesWritten()))
		if ww.Status() >= 500 {
			RequestErrors.WithLabelValues(r.URL.Path, r.Method, status).Inc()
		}
	}
	return http.HandlerFunc(fn)
}
