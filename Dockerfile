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
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/app/mockhub"]
