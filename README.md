# MockHub

MockHub is a local mock server for DataHub APIs. It serves canned responses for GraphQL and REST endpoints and allows toggling scenarios via request headers.

## Running

```bash
cp .env.example .env
make run
```

GraphQL endpoints listen on `:8082`, REST endpoints on `:8083`, and the docs/UI on `:3002` by default.

### Docker

```bash
make docker-build
make docker-run
```

Environment variables can be passed with `-e`, for example `docker run -e GRAPHQL_PORT=8082 -e REST_PORT=8083 -e UI_PORT=3002 -p 8082:8082 -p 8083:8083 -p 3002:3002 mockhub`.

- GraphQL endpoint: `POST http://localhost:8082/api/graphql`
- GraphiQL UI: `http://localhost:3002/graphiql`
- REST endpoints: `http://localhost:8083/api/...`
- Swagger UI: `http://localhost:3002/swagger`
- Example frontend: `http://localhost:3002/`
- Server logs: `http://localhost:3002/logs` (also shown in the frontend)

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
