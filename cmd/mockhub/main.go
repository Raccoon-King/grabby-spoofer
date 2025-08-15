// Copyright (c) 2024. Licensed under the MIT License.
package main

import (
	"log"

	"github.com/example/mockhub/internal/graphql"
	"github.com/example/mockhub/internal/httpserver"
	"github.com/example/mockhub/internal/rest"
)

func main() {
	if err := graphql.LoadMocks("internal/mocks/graphql"); err != nil {
		log.Fatalf("load graphql mocks: %v", err)
	}
	if err := rest.LoadMocks("internal/mocks/rest"); err != nil {
		log.Fatalf("load rest mocks: %v", err)
	}
	if err := httpserver.Run(); err != nil {
		log.Fatal(err)
	}
}
