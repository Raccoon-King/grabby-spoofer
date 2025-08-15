// Copyright (c) 2024. Licensed under the MIT License.
package graphql

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlerSearchDefault(t *testing.T) {
	if err := LoadMocks("../mocks/graphql"); err != nil {
		t.Fatal(err)
	}
	reqBody := []byte(`{"query":"query search{search{start}}","operationName":"search"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/graphql", bytes.NewReader(reqBody))
	rr := httptest.NewRecorder()
	Handler(rr, req)
	expected, _ := os.ReadFile("../mocks/graphql/default/search.json")
	if rr.Code != http.StatusOK {
		t.Fatalf("code %d", rr.Code)
	}
	if !bytes.Equal(bytes.TrimSpace(rr.Body.Bytes()), bytes.TrimSpace(expected)) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestHandlerSearchError(t *testing.T) {
	if err := LoadMocks("../mocks/graphql"); err != nil {
		t.Fatal(err)
	}
	reqBody := []byte(`{"query":"query search{search{start}}","operationName":"search"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/graphql", bytes.NewReader(reqBody))
	req.Header.Set("X-Mock-Scenario", "error")
	rr := httptest.NewRecorder()
	Handler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("code %d", rr.Code)
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte("INTERNAL_ERROR")) {
		t.Fatalf("expected error, got %s", rr.Body.String())
	}
}
