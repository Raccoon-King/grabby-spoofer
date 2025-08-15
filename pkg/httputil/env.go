// Copyright (c) 2024. Licensed under the MIT License.
package httputil

import (
	"net/http"
	"os"
)

// GetEnv reads an env var or returns default.
func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ScenarioFromHeader returns the mock scenario from request header.
func ScenarioFromHeader(r *http.Request) string {
	s := r.Header.Get("X-Mock-Scenario")
	if s == "" {
		return "default"
	}
	return s
}
