// Copyright (c) 2024. Licensed under the MIT License.
package httputil

import (
	"encoding/json"
	"net/http"
)

// JSON writes value as JSON with status.
func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Error writes an error message as JSON.
func Error(w http.ResponseWriter, status int, code, msg string) {
	JSON(w, status, map[string]string{"code": code, "message": msg})
}
