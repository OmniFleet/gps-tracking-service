package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock HealthCheck Models
type MockChecker struct {
	Error error
	Data  map[string]string
}

func (m MockChecker) Alive() (map[string]string, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	if m.Data != nil {
		return map[string]string{
			"alive": "true",
		}, nil
	}
	return m.Data, nil
}

func (m MockChecker) Ready() (map[string]string, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	if m.Data != nil {
		return map[string]string{
			"ready": "true",
		}, nil
	}
	return m.Data, nil
}

func TestLivenessHealthCheck(t *testing.T) {

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "LivenessHealthy",
			status: http.StatusOK,
		},
		{
			name:   "LivenessError",
			status: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var mock MockChecker
			if tt.status > 499 {
				mock.Error = fmt.Errorf("I should fail")
			}

			Liveness(mock).ServeHTTP(rr, r)
			rs := rr.Result()
			if rs.StatusCode != tt.status {
				t.Errorf("got %d; expected %d", rs.StatusCode, tt.status)
			}
		})
	}
}

func TestReadyHealthCheck(t *testing.T) {

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "ReadyHealthy",
			status: http.StatusOK,
		},
		{
			name:   "ReadyError",
			status: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var mock MockChecker
			if tt.status > 499 {
				mock.Error = fmt.Errorf("I should fail")
			}

			Ready(mock).ServeHTTP(rr, r)
			rs := rr.Result()
			if rs.StatusCode != tt.status {
				t.Errorf("got %d; expected %d", rs.StatusCode, tt.status)
			}
		})
	}
}
