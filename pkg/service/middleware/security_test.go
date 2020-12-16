package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "X-Frame-Options",
			want: "deny",
		},
		{
			name: "X-XSS-Protection",
			want: "1; mode=block",
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
				w.Write([]byte("OK"))
			})

			SecurityHeaders(next).ServeHTTP(rr, r)

			rs := rr.Result()

			header := rs.Header.Get(tt.name)
			if header != tt.want {
				t.Errorf("want %q; got %q", tt.want, header)
			}

			if rs.StatusCode != http.StatusOK {
				t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
			}
			defer rs.Body.Close()

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

		SecurityHeaders(next).ServeHTTP(rr, r)

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
