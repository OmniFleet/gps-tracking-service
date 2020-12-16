package service

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"scbunn.org/tmp/gps-tracking-service/pkg/service/handlers"
	"scbunn.org/tmp/gps-tracking-service/pkg/service/middleware"
)

func (s *Service) Routes() http.Handler {
	// configure a sampled logger for our routes
	sampledLogger := s.logger.Sample(&zerolog.BurstSampler{
		Burst:       15,
		Period:      1 * time.Second,
		NextSampler: &zerolog.BasicSampler{N: 100},
	})

	middleware.RegisterMetrics()

	r := chi.NewRouter()
	r.Use(middleware.PrometheusTelemetry)
	r.Use(mw.RequestID)
	r.Use(middleware.ZeroLog(&sampledLogger))
	r.Use(mw.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.SecurityHeaders)
	r.Use(mw.Recoverer)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/health", func(r chi.Router) {
		r.Get("/liveness", handlers.Liveness(s.telemetry))
		r.Get("/readiness", handlers.Ready(s.telemetry))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/location/{id}", handlers.GetLocation(s.telemetry))
		r.Post("/location/", handlers.UpdateLocation(s.telemetry))
		r.Get("/location/", handlers.GetAllLocations(s.telemetry))
	})

	return r
}
