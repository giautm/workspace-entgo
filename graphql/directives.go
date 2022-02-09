package graphql

import (
	"context"

	"giautm.dev/awesome/graphql/generated"
	"giautm.dev/awesome/graphql/model"
	"giautm.dev/awesome/internal/auth"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	errNotAuthorized = gqlerror.Errorf("you are not authorized for this resource")
)

func directiveAuth(ctx context.Context, obj interface{}, next graphql.Resolver, requires *model.Role) (res interface{}, err error) {
	token := auth.TokenFromContext(ctx)
	if token == nil {
		return nil, errNotAuthorized
	}

	if requires != nil && *requires == model.RoleAdmin {
		if c, ok := token.Claims.(*auth.LegacyClaims); ok {
			if c.Guard == auth.GuardEmployer {
				return next(ctx)
			}
		}

		return nil, errNotAuthorized
	}

	return next(ctx)
}

var directives = generated.DirectiveRoot{
	Auth: directiveAuth,
}
