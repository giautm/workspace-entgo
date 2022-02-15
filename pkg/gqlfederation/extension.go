package gqlfederation

import (
	"embed"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

const GQLFederationAnnotation = "GQLFederation"

var (
	GQLFederationTemplate = parseT("template/gql_federation.tmpl")

	AllTemplates = []*gen.Template{
		GQLFederationTemplate,
	}

	//go:embed template/*
	templates embed.FS
)

func parseT(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(entgql.TemplateFuncs).
		ParseFS(templates, path))
}

type (
	Config struct {
		FederatedNodes []string
	}

	Extension struct {
		entc.DefaultExtension

		templates []*gen.Template
		hooks     []gen.Hook
		cfg       *Config
	}

	// ExtensionOption allows for managing the Extension configuration
	// using functional options.
	ExtensionOption func(*Extension) error
)

var (
	_ entc.Extension = (*Extension)(nil)
)

func WithFederatedNodes(nodes ...string) ExtensionOption {
	return func(e *Extension) error {
		if e.cfg == nil {
			e.cfg = &Config{}
		}

		e.cfg.FederatedNodes = nodes
		return nil
	}
}

func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{templates: AllTemplates}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}

	ex.hooks = append(ex.hooks, func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			if ex.cfg == nil {
				return next.Generate(g)
			}

			if g.Annotations == nil {
				g.Annotations = gen.Annotations{}
			}

			g.Annotations[GQLFederationAnnotation] = ex.cfg
			return next.Generate(g)
		})
	})

	return ex, nil
}

// Hooks of the extension.
func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

// Templates of the extension.
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}
