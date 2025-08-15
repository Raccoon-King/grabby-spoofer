# Copyright (c) 2024. Licensed under the MIT License.
# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o grabby-api ./cmd/grabby-api

FROM alpine:3.19
WORKDIR /app
COPY --from=build /src/grabby-api ./grabby-api
COPY --from=build /src/examples/frontend ./examples/frontend
COPY --from=build /src/internal/rest/openapi.yaml ./internal/rest/openapi.yaml
EXPOSE 3002 8082
ENV API_PORT=3002
ENV UI_PORT=8082
ENTRYPOINT ["/app/grabby-api"]
