package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

type MockModel struct {
	Error      error
	GetAllSize int
}

func (m MockModel) NewTelemetry() *models.Telemetry {
	return &models.Telemetry{
		Source:   "ci test data",
		ObjectID: "0001",
		Status:   "TESTING",
		Position: models.Position{
			Latitude:  127.000,
			Longitude: -127.000,
			Elevation: 100,
		},
	}
}

func (m MockModel) Add(t models.Telemetry) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	return t.Id, nil
}

func (m MockModel) Get(id string) (*models.Telemetry, error) {
	if m.Error != nil {
		return nil, models.ErrNoRecord
	}
	t := m.NewTelemetry()
	t.Id = id
	return t, nil
}

func (m MockModel) GetAll() []models.Telemetry {
	var results []models.Telemetry
	for i := 0; i < m.GetAllSize; i++ {
		id := uuid.New().String()
		t := m.NewTelemetry()
		t.Id = id
		results = append(results, *t)
	}
	return results
}

func TestGetLocation(t *testing.T) {
	tests := []struct {
		name   string
		mock   MockModel
		status int
	}{
		{name: "ValidId", mock: MockModel{}, status: http.StatusOK},
		{name: "IdNotFound", mock: MockModel{Error: models.ErrNoRecord}, status: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			id := uuid.New().String()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", id)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			GetLocation(tt.mock).ServeHTTP(w, r)
			rs := w.Result()
			defer rs.Body.Close()
			if rs.StatusCode != tt.status {
				t.Errorf("expected %d status; got %d status", tt.status, rs.StatusCode)
			}
		})
	}

}

func TestGetAllLocations(t *testing.T) {
	tests := []struct {
		name   string
		mock   MockModel
		status int
	}{
		{name: "EmptyDatabase", mock: MockModel{}, status: http.StatusOK},
		{name: "OneItem", mock: MockModel{GetAllSize: 1}, status: http.StatusOK},
		{name: "HundredItems", mock: MockModel{GetAllSize: 100}, status: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			GetAllLocations(tt.mock).ServeHTTP(w, r)
			rs := w.Result()
			defer rs.Body.Close()

			if rs.StatusCode != tt.status {
				t.Errorf("expected %d status; got %d status", tt.status, rs.StatusCode)
			}
		})
	}

}

func TestUpdateLocation(t *testing.T) {
	tests := []struct {
		name   string
		mock   MockModel
		body   string
		status int
	}{
		{name: "InvalidJSON", mock: MockModel{}, body: `{"foo": "bar"}`, status: http.StatusBadRequest},
		{name: "ValidRequest", mock: MockModel{}, body: `{"source": "testing", "objectId": "123", "position": {"latitude": 123, "longitude": -123}}`, status: http.StatusCreated},
		{name: "InternalError", mock: MockModel{Error: fmt.Errorf("bad thing")}, body: `{"source": "testing", "objectId": "123", "position": {"latitude": 123, "longitude": -123}}`, status: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.NewReader(tt.body)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", b)

			UpdateLocation(tt.mock).ServeHTTP(w, r)
			rs := w.Result()
			defer rs.Body.Close()
			body, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				t.Errorf("%s", err.Error())
			}
			t.Logf("body: %s\n", body)

			if rs.StatusCode != tt.status {
				t.Errorf("expected %d status; got %d status", tt.status, rs.StatusCode)
			}
		})
	}
}
