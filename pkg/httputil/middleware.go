// Copyright (c) 2024. Licensed under the MIT License.
package httputil

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// RequestLogger logs basic request information.
func RequestLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()
			scenario := ScenarioFromHeader(r)
			ww := &responseWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(ww, r)
			logger.Info().Str("method", r.Method).Str("path", r.URL.Path).
				Int("status", ww.status).Dur("latency", time.Since(start)).
				Str("scenario", scenario).Msg("request")
		})
	}
}

// Recoverer recovers from panics.
func Recoverer(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error().Interface("err", err).Msg("panic")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
