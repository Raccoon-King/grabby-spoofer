#!/bin/bash
# Copyright (c) 2024. Licensed under the MIT License.
set -euo pipefail

GRAPHQL_URL=${GRAPHQL_URL:-https://demo.datahubproject.io/api/graphql}
OPENAPI_URL=${OPENAPI_URL:-https://demo.datahubproject.io/openapi/openapi.yaml}
AUTH_TOKEN=${AUTH_TOKEN:-}

if [[ -n "$AUTH_TOKEN" ]]; then
  AUTH_HEADER="Authorization: Bearer $AUTH_TOKEN"
else
  AUTH_HEADER=""
fi

# Fetch GraphQL introspection
curl -sSf ${AUTH_HEADER:+-H "$AUTH_HEADER"} \
  -H 'Content-Type: application/json' \
  -d @- "$GRAPHQL_URL" <<'QRY' > internal/mocks/graphql/default/introspection.json
{"query":"query IntrospectionQuery { __schema { queryType { name } mutationType { name } subscriptionType { name } types { ...FullType } directives { name description locations args { ...InputValue } } } } fragment FullType on __Type { kind name description fields(includeDeprecated: true) { name description args { ...InputValue } type { ...TypeRef } isDeprecated deprecationReason } inputFields { ...InputValue } interfaces { ...TypeRef } enumValues(includeDeprecated: true) { name description isDeprecated deprecationReason } possibleTypes { ...TypeRef } } fragment InputValue on __InputValue { name description type { ...TypeRef } defaultValue } fragment TypeRef on __Type { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name } } } } } } }"}
QRY

# Fetch OpenAPI spec
curl -sSf ${AUTH_HEADER:+-H "$AUTH_HEADER"} "$OPENAPI_URL" > internal/rest/openapi.yaml

echo "Definitions updated"
