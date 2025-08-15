// Copyright (c) 2024. Licensed under the MIT License.
package rest

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/example/grabby-api/pkg/httputil"
)

var restMocks = map[string]map[string]map[string]interface{}{}

// LoadMocks loads REST mock JSON files.
func LoadMocks(base string) error {
	endpoints := []string{"search", "entity", "tags_create", "entity_tags", "entity_tag_delete", "ownership_get", "ownership_post", "ownership_delete"}
	scenarios := []string{"default", "empty", "partial", "error", "auth"}
	for _, ep := range endpoints {
		restMocks[ep] = map[string]map[string]interface{}{}
		for _, sc := range scenarios {
			path := filepath.Join(base, sc, ep+".json")
			b, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var v map[string]interface{}
			_ = json.Unmarshal(b, &v)
			restMocks[ep][sc] = v
		}
	}
	return nil
}

// Register registers REST routes under /api prefix.
func Register(r chi.Router) {
	r.Get("/search", searchHandler)
	r.Get("/entities/{urn}", entityHandler)
	r.Post("/tags", tagCreateHandler)
	r.Post("/entities/{urn}/tags", entityTagAttachHandler)
	r.Delete("/entities/{urn}/tags/{tag}", entityTagDeleteHandler)
	r.Get("/datasets/{urn}/ownership", ownershipGetHandler)
	r.Post("/datasets/{urn}/ownership", ownershipPostHandler)
	r.Delete("/datasets/{urn}/ownership/{ownerUrn}", ownershipDeleteHandler)
}

func scenario(r *http.Request) string { return httputil.ScenarioFromHeader(r) }

func respond(w http.ResponseWriter, r *http.Request, ep string) {
	sc := scenario(r)
	data := restMocks[ep][sc]
	if data == nil {
		data = restMocks[ep]["default"]
	}
	if sc == "error" {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "mock error")
		return
	}
	if sc == "auth" {
		w.Header().Set("WWW-Authenticate", "Mock")
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Auth required")
		return
	}
	httputil.JSON(w, http.StatusOK, data)
}

func searchHandler(w http.ResponseWriter, r *http.Request)          { respond(w, r, "search") }
func entityHandler(w http.ResponseWriter, r *http.Request)          { respond(w, r, "entity") }
func tagCreateHandler(w http.ResponseWriter, r *http.Request)       { respond(w, r, "tags_create") }
func entityTagAttachHandler(w http.ResponseWriter, r *http.Request) { respond(w, r, "entity_tags") }
func entityTagDeleteHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, "entity_tag_delete")
}
func ownershipGetHandler(w http.ResponseWriter, r *http.Request)  { respond(w, r, "ownership_get") }
func ownershipPostHandler(w http.ResponseWriter, r *http.Request) { respond(w, r, "ownership_post") }
func ownershipDeleteHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, "ownership_delete")
}

// SwaggerUI serves Swagger UI
func SwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html><html><head><title>Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" /></head>
    <body><div id="swagger"></div>
    <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
    <script>SwaggerUIBundle({url:'/openapi.yaml',dom_id:'#swagger'});</script></body></html>`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

// OpenAPI serves the openapi spec.
func OpenAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := os.ReadFile("internal/rest/openapi.yaml")
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
