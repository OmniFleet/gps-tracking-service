package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

func TestRegisteredMetrics(t *testing.T) {
	// no panic, good to go
	RegisterMetrics()
}

func TestRequestDurationMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		errCount int
	}{
		{
			name:     "Valid Request",
			status:   http.StatusOK,
			errCount: 0,
		},
		{
			name:     "Internal Server Error",
			status:   http.StatusInternalServerError,
			errCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a mock HTTP handler that we can pass to the SecurityHeaders
			// middleware
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte("OK"))
			})

			PrometheusTelemetry(next).ServeHTTP(rr, r)

			rs := rr.Result()
			defer rs.Body.Close()

			if rs.StatusCode != tt.status {
				t.Errorf("next sent %q; got %q", tt.status, rs.StatusCode)
			}

			labels := []string{"/", http.MethodGet, http.StatusText(tt.status)}
			errorCount, _ := models.GetCounterValue(RequestErrors, labels...)
			if errorCount > float64(tt.errCount) {
				t.Errorf("got %f errors; expected %d errors", errorCount, tt.errCount)
			}

		})
	}

	t.Run("CallNextHandler", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP handler that we can pass to the SecurityHeaders
		// middleware
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		PrometheusTelemetry(next).ServeHTTP(rr, r)

		rs := rr.Result()

		if rs.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
		}
		defer rs.Body.Close()

		// check if the next handler is called
		body, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != "OK" {
			t.Errorf("next handler returns OK; got %s", string(body))
		}

	})

}
