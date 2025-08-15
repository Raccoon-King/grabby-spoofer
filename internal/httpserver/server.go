// Copyright (c) 2024. Licensed under the MIT License.
package httpserver

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/example/mockhub/internal/graphql"
	"github.com/example/mockhub/internal/rest"
	"github.com/example/mockhub/pkg/httputil"
)

// NewAPIServer configures the backend API server.
func NewAPIServer() *http.Server {
	port := httputil.GetEnv("API_PORT", "3002")
	logLevel := httputil.GetEnv("LOG_LEVEL", "debug")
	allowed := httputil.GetEnv("CORS_ALLOWED_ORIGINS", "*")

	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	cw := zerolog.ConsoleWriter{Out: io.MultiWriter(os.Stderr, httputil.GlobalLogBuffer)}
	logger := log.Output(cw)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowed},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))
	r.Use(httputil.RequestLogger(logger))
	r.Use(httputil.Recoverer(logger))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.Route("/api", func(api chi.Router) {
		api.Post("/graphql", graphql.Handler)
		api.Get("/graphql", graphql.Handler)
		rest.Register(api)
	})

	return &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

// NewUIServer configures the frontend/docs server.
func NewUIServer() *http.Server {
	port := httputil.GetEnv("UI_PORT", "8082")
	logLevel := httputil.GetEnv("LOG_LEVEL", "debug")

	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	cw := zerolog.ConsoleWriter{Out: io.MultiWriter(os.Stderr, httputil.GlobalLogBuffer)}
	logger := log.Output(cw)

	r := chi.NewRouter()
	r.Use(httputil.RequestLogger(logger))
	r.Use(httputil.Recoverer(logger))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/logs", httputil.LogsHandler)

	r.Get("/graphiql", graphql.GraphiQL)
	r.Get("/swagger", rest.SwaggerUI)
	r.Get("/openapi.yaml", rest.OpenAPI)
	r.Handle("/*", http.FileServer(http.Dir("examples/frontend")))

	return &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

// Run starts both servers and blocks until one exits.
func Run() error {
	api := NewAPIServer()
	ui := NewUIServer()

	errCh := make(chan error, 2)
	go func() { errCh <- api.ListenAndServe() }()
	go func() { errCh <- ui.ListenAndServe() }()
	return <-errCh
}
