package graphql

import (
	"go.uber.org/fx"

	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/graphql/resolver"
	"giautm.dev/awesome/internal/project"
)

var Module = fx.Options(
	fx.Provide(NewServeFx),
)

func NewServeFx(client *ent.Client) (*Server, error) {
	opts := ProductionOptions
	if project.DevMode() {
		opts = DevelopmentOptions
	}

	opts = append(opts, WithEntTransaction(client))
	return NewServer(resolver.NewResolver(client), opts...)
}
