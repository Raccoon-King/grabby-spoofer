# Examples

```bash
# GraphQL search
curl -H 'Content-Type: application/json' -d '{"query":"query search{search{start}}","operationName":"search"}' http://localhost:3002/api/graphql

# REST search
curl http://localhost:3002/api/search

# Error scenario
curl -H 'X-Mock-Scenario: error' http://localhost:3002/api/search
```
