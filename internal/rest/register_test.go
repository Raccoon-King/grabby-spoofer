// Copyright (c) 2024. Licensed under the MIT License.
package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setup() http.Handler {
	if err := LoadMocks("../mocks/rest"); err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	Register(r)
	return r
}

func TestSearchDefault(t *testing.T) {
	h := setup()
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Sample Dataset") {
		t.Fatalf("unexpected body %s", rr.Body.String())
	}
}

func TestSearchAuth(t *testing.T) {
	h := setup()
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	req.Header.Set("X-Mock-Scenario", "auth")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}
