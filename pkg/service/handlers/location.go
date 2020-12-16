package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

func GetLocation(t models.TelemetryReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		location, err := t.Get(id)
		if err != nil {
			renderError(w, http.StatusNotFound, err)
			return
		}
		renderJSON(w, http.StatusOK, location)
	}
}

func GetAllLocations(t models.TelemetryReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locations := t.GetAll()
		if len(locations) == 0 {
			renderJSON(w, http.StatusOK, []string{})
			return
		}
		renderJSON(w, http.StatusOK, locations)
	}
}

func UpdateLocation(t models.TelemetryReaderWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		var telemetry models.Telemetry

		if err := telemetry.FromJSON(r); err != nil {
			renderError(w, http.StatusBadRequest, err)
			return
		}
		telemetry.Updated = now
		telemetry.Id = fmt.Sprintf("%s-%s", telemetry.Source, telemetry.ObjectID)

		id, err := t.Add(telemetry)
		if err != nil {
			renderError(w, http.StatusInternalServerError, err)
			return
		}

		renderJSON(w, http.StatusCreated, map[string]string{"message": "created", "id": id})
	}
}

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(response); err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}
}

func renderError(w http.ResponseWriter, status int, err error) {
	data := map[string]string{
		"error":   err.Error(),
		"status":  fmt.Sprintf("%d", status),
		"message": "an error has occured",
	}

	response, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(response); err != nil {
		panic(err)
	}
}
