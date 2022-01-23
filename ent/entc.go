//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	exEntGQL, err := entgql.NewExtension(
		entgql.WithWhereFilters(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		// Generate the filters to a separate schema
		// file and load it in the gqlgen.yml config.
		entgql.WithSchemaPath("./internal/graphql/schema/ent.gql"),
	)
	if err != nil {
		log.Fatalf("entc: creating EntGQL extension: %v", err)
	}

	err = entc.Generate("./ent/schema", &gen.Config{},
		entc.Extensions(exEntGQL),
		entc.TemplateDir("./ent/template"),
		entc.TemplateDir("./ent/schema/pulid/template"),
	)
	if err != nil {
		log.Fatalf("entc: running ent codegen: %v", err)
	}
}
