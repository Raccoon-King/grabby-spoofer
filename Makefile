# Copyright (c) 2024. Licensed under the MIT License.
SHELL := /bin/bash

.PHONY: build run test docker-build docker-run import-defs

build:
	go build ./cmd/mockhub

run:
	go run ./cmd/mockhub

test:
        go test ./...

docker-build:
        docker build -t mockhub .

docker-run: docker-build
        docker run --rm -p 3002:3002 -p 8082:8082 mockhub

import-defs:
        scripts/import-defs.sh
