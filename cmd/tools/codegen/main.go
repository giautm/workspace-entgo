package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"giautm.dev/awesome/pkg/gqlfederation"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
)

const gqlgenConfigFile = "./gqlgen.yml"

func main() {
	gqlgenOpts := []api.Option{}

	exEntGQL, err := entgql.NewExtension(
		entgql.WithWhereFilters(true),
		entgql.WithConfigPath(gqlgenConfigFile, gqlgenOpts...),
		// Generate the filters to a separate schema
		// file and load it in the gqlgen.yml config.
		entgql.WithSchemaPath("./graphql/schema/ent.gql"),
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

	cfg, err := config.LoadConfig(gqlgenConfigFile)
	if err != nil {
		log.Fatalf("gqlgen: failed to load config: %v", err)
	}

	if err = api.Generate(cfg, gqlgenOpts...); err != nil {
		log.Fatalf("gqlgen: running generate: %v", err)
	}
}
