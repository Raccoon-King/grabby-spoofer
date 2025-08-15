# Usage Guide

This guide shows how to interact with MockHub's mocked DataHub APIs.

## Prerequisites
* Go 1.21+
* Server running: `make run`

## Switching Scenarios
Responses are controlled by the `X-Mock-Scenario` header. Supported values:
`default`, `empty`, `partial`, `error`, `auth`. Omit the header for `default`.

## GraphQL
### Search
```
curl -H 'Content-Type: application/json' \
  -d '{"query":"query search{search{start total}}","operationName":"search"}' \
  http://localhost:8080/api/graphql
```

### Entity Lookup with Scenario
```
curl -H 'Content-Type: application/json' \
  -H 'X-Mock-Scenario: error' \
  -d '{"query":"query entityByUrn($urn:String!){entityByUrn(urn:$urn){urn}}","variables":{"urn":"urn:li:dataset:1"},"operationName":"entityByUrn"}' \
  http://localhost:8080/api/graphql
```

## REST
### Search
```
curl http://localhost:8080/api/search?query=sample
```

### Auth Scenario
```
curl -i -H 'X-Mock-Scenario: auth' http://localhost:8080/api/search
```

## Frontend
A small HTML frontend is provided under `examples/frontend`. Open
`examples/frontend/index.html` in a browser and make requests without using curl.
