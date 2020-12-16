package handlers

import (
	"net/http"

	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

func Liveness(t models.HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := t.Alive()
		if err != nil {
			renderError(w, http.StatusServiceUnavailable, err)
			return
		}
		renderJSON(w, http.StatusOK, data)
	}
}

func Ready(t models.HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := t.Ready()
		if err != nil {
			renderError(w, http.StatusServiceUnavailable, err)
			return
		}
		renderJSON(w, http.StatusOK, data)
	}
}
