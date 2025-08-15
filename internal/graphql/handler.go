// Copyright (c) 2024. Licensed under the MIT License.
package graphql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/example/mockhub/pkg/httputil"
)

// request represents incoming GraphQL request.
type request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

var mockData = map[string]map[string]json.RawMessage{}

// LoadMocks loads JSON files from internal/mocks/graphql.
func LoadMocks(base string) error {
	ops := []string{"searchAcrossEntities", "search", "entityByUrn"}
	scenarios := []string{"default", "empty", "partial", "error", "auth"}
	for _, op := range ops {
		mockData[op] = map[string]json.RawMessage{}
		for _, sc := range scenarios {
			path := filepath.Join(base, sc, op+".json")
			b, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}
			mockData[op][sc] = b
		}
	}
	// introspection
	b, err := ioutil.ReadFile(filepath.Join(base, "default", "introspection.json"))
	if err == nil {
		mockData["IntrospectionQuery"] = map[string]json.RawMessage{"default": b}
	}
	return nil
}

// Handler serves GraphQL requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// simple health for GET
		httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	op := req.OperationName
	if op == "" {
		// try detect by query substring
		if strings.Contains(req.Query, "__schema") {
			op = "IntrospectionQuery"
		} else if strings.Contains(req.Query, "searchAcrossEntities") {
			op = "searchAcrossEntities"
		} else if strings.Contains(req.Query, "search(") {
			op = "search"
		} else if strings.Contains(req.Query, "browse(") {
			op = "browse"
		} else if strings.Contains(req.Query, "entityByUrn") {
			op = "entityByUrn"
		}
	}
	scenario := httputil.ScenarioFromHeader(r)
	b := mockData[op][scenario]
	if len(b) == 0 {
		// fallback to default
		b = mockData[op]["default"]
	}
	if len(b) == 0 {
		httputil.Error(w, http.StatusNotFound, "NOT_FOUND", "no mock for operation")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GraphiQL serves a tiny GraphiQL page.
func GraphiQL(w http.ResponseWriter, r *http.Request) {
	endpoint := fmt.Sprintf("http://localhost:%s/api/graphql", httputil.GetEnv("API_PORT", "3002"))
	html := fmt.Sprintf(`<!DOCTYPE html><html><head><title>GraphiQL</title>
    <link href="https://unpkg.com/graphiql/graphiql.min.css" rel="stylesheet" />
    </head><body style=\"margin:0;\">
    <div id=\"graphiql\" style=\"height:100vh;\"></div>
    <script crossorigin src=\"https://unpkg.com/react/umd/react.production.min.js\"></script>
    <script crossorigin src=\"https://unpkg.com/react-dom/umd/react-dom.production.min.js\"></script>
    <script src=\"https://unpkg.com/graphiql/graphiql.min.js\"></script>
    <script>const graphQLFetcher = graphQLParams => fetch('%s',{method:'post',headers:{'Content-Type':'application/json'},body:JSON.stringify(graphQLParams)}).then(r=>r.json());ReactDOM.render(React.createElement(GraphiQL,{fetcher:graphQLFetcher}),document.getElementById('graphiql'));</script>
    </body></html>`, endpoint)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}
