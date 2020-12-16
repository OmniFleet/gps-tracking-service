package middleware

import (
	"net/http"

	mw "github.com/go-chi/chi/middleware"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// get the current request id
		id := mw.GetReqID(r.Context())
		w.Header().Set("X-Request-Id", id)

		next.ServeHTTP(w, r)
	})
}
