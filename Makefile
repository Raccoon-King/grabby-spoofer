# Copyright (c) 2024. Licensed under the MIT License.
SHELL := /bin/bash

.PHONY: build run test docker-build docker-run import-defs

build:
go build ./cmd/grabby-api

run:
go run ./cmd/grabby-api

test:
        go test ./...

docker-build:
docker build -t grabby-api .

docker-run: docker-build
docker run --rm -p 3002:3002 -p 8082:8082 grabby-api

import-defs:
        scripts/import-defs.sh
