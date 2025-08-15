// Copyright (c) 2024. Licensed under the MIT License.
package httpserver

import (
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

// New creates a configured http.Server.
func New() *http.Server {
	// Configuration via env vars
	port := httputil.GetEnv("PORT", "8080")
	logLevel := httputil.GetEnv("LOG_LEVEL", "info")
	allowed := httputil.GetEnv("CORS_ALLOWED_ORIGINS", "*")

	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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

	// GraphQL
	r.Route("/api", func(api chi.Router) {
		api.Post("/graphql", graphql.Handler)
		api.Get("/graphql", graphql.Handler) // for graphiql / health
		// REST subset
		rest.Register(api)
	})

	// GraphiQL UI
	r.Get("/graphiql", graphql.GraphiQL)
	// Swagger UI
	r.Get("/swagger", rest.SwaggerUI)
	r.Get("/openapi.yaml", rest.OpenAPI)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return server
}

// Run starts the server and blocks.
func Run() error {
	server := New()
	return server.ListenAndServe()
}
