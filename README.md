# Grabby API

Grabby API is a local mock server for DataHub APIs. It serves canned responses for GraphQL and REST endpoints and allows toggling scenarios via request headers.

## Running

```bash
cp .env.example .env
make run
```

API server listens on `:3002` and the UI on `:8082` by default.

### Docker

```bash
make docker-build
make docker-run
```

Environment variables can be passed with `-e`, for example `docker run -e API_PORT=3002 -e UI_PORT=8082 -p 3002:3002 -p 8082:8082 grabby-api`.

- GraphQL endpoint: `POST http://localhost:3002/api/graphql`
- GraphiQL UI: `http://localhost:8082/graphiql`
- REST endpoints: `http://localhost:3002/api/...`
- Swagger UI: `http://localhost:8082/swagger`
- Example frontend: `http://localhost:8082/`

Detailed usage instructions are available in [USAGE.md](USAGE.md).

Scenarios are selected with header `X-Mock-Scenario` and can be `default`, `empty`, `partial`, `error`, or `auth`.

## Development

- `make build`
- `make test`

Mock payloads live under `internal/mocks`.

### Updating API definitions

If upstream DataHub APIs change, refresh the GraphQL introspection file and OpenAPI spec with:

```bash
AUTH_TOKEN=yourtoken make import-defs
```

Customize the source URLs via `GRAPHQL_URL` and `OPENAPI_URL` environment variables.
