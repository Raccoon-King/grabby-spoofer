# MockHub

MockHub is a local mock server for DataHub APIs. It serves canned responses for GraphQL and REST endpoints and allows toggling scenarios via request headers.

## Running

```bash
cp .env.example .env
make run
```

Server listens on `:8080` by default.

### Docker

```bash
make docker-build
make docker-run
```

Environment variables can be passed with `-e`, for example `docker run -e PORT=8080 -p 8080:8080 mockhub`.

- GraphQL endpoint: `POST /api/graphql`
- GraphiQL UI: `/graphiql`
- REST endpoints: `/api/...`
- Swagger UI: `/swagger`
- Example frontend: open `examples/frontend/index.html`

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
