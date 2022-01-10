//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatalf("gqlgen: failed to load config: %v", err)
	}

	if err = api.Generate(cfg); err != nil {
		log.Fatalf("gqlgen: running generate: %v", err)
	}
}
