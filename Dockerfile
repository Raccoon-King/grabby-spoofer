# Copyright (c) 2024. Licensed under the MIT License.
# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mockhub ./cmd/mockhub

FROM alpine:3.19
WORKDIR /app
COPY --from=build /src/mockhub ./mockhub
COPY --from=build /src/examples/frontend ./examples/frontend
COPY --from=build /src/internal/rest/openapi.yaml ./internal/rest/openapi.yaml
EXPOSE 3002 8082 8083
ENV GRAPHQL_PORT=8082
ENV REST_PORT=8083
ENV UI_PORT=3002
ENTRYPOINT ["/app/mockhub"]
