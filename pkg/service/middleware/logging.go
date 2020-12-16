package middleware

import (
	"net/http"
	"time"

	mw "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

func ZeroLog(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()

			ww := mw.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			defer func() {
				end := time.Now()
				log.Info().
					Timestamp().
					Dur("duration", end.Sub(start)).
					Fields(map[string]interface{}{
						"remoteIP":  r.RemoteAddr,
						"url":       r.URL.Path,
						"proto":     r.Proto,
						"method":    r.Method,
						"userAgent": r.Header.Get("User-Agent"),
						"status":    ww.Status(),
						"bytesIn":   r.ContentLength,
						"bytesOut":  ww.BytesWritten(),
						"requestId": mw.GetReqID(r.Context()),
					}).
					Msg("request recieved")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
