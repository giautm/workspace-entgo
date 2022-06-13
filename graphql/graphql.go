package graphql

import (
	"go.uber.org/fx"

	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/graphql/resolver"
	"giautm.dev/awesome/internal/project"
)

// Module exports the graphql module.
var Module = fx.Options(
	fx.Provide(NewServeFx),
)

// NewServeFx returns a new graphql server.
func NewServeFx(client *ent.Client) (*Server, error) {
	opts := ProductionOptions
	if project.DevMode() {
		opts = DevelopmentOptions
	}

	return NewServer(resolver.NewResolver(), client, opts...)
}
