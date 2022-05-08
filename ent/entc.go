//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"giautm.dev/awesome/pkg/gqlfederation"
)

func main() {
	exEntGQL, err := entgql.NewExtension(
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("./graphql/schema/ent.gql"),
		entgql.WithWhereFilters(true),
	)
	if err != nil {
		log.Fatalf("entc: creating EntGQL extension: %v", err)
	}

	exFederation, err := gqlfederation.NewExtension()
	if err != nil {
		log.Fatalf("entc: creating GQLFederation extension: %v", err)
	}
	err = entc.Generate("./ent/schema", &gen.Config{},
		entc.Extensions(exEntGQL, exFederation),
		entc.TemplateDir("./ent/template"),
	)
	if err != nil {
		log.Fatalf("entc: running ent codegen: %v", err)
	}
}
