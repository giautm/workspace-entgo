package graphql

import (
	"context"

	"giautm.dev/awesome/graphql/generated"
	"giautm.dev/awesome/graphql/model"
	"giautm.dev/awesome/internal/auth"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func directiveAuth(ctx context.Context, obj interface{}, next graphql.Resolver, requires *model.Role) (res interface{}, err error) {
	claims := auth.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, gqlerror.Errorf("you are not authorized for this resource")
	}

	if requires != nil && *requires == model.RoleAdmin {
		if claims.Guard == auth.GuardEmployer {
			return next(ctx)
		}

		return nil, gqlerror.Errorf("you are not authorized for this resource")
	}

	return next(ctx)
}

var directives = generated.DirectiveRoot{
	Auth: directiveAuth,
}
