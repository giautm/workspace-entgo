package directive

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"giautm.dev/awesome/internal/auth"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type AuthRole string

const (
	AuthRoleAdmin AuthRole = "ADMIN"
)

var AllRole = []AuthRole{
	AuthRoleAdmin,
}

func (e AuthRole) IsValid() bool {
	switch e {
	case AuthRoleAdmin:
		return true
	}
	return false
}

func (e AuthRole) String() string {
	return string(e)
}

func (e *AuthRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AuthRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e AuthRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver, requires *AuthRole) (res interface{}, err error) {
	claims := auth.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, gqlerror.Errorf("you are not authorized for this resource")
	}

	if requires != nil && *requires == AuthRoleAdmin {
		if claims.Guard == auth.GuardEmployer {
			return next(ctx)
		}

		return nil, gqlerror.Errorf("you are not authorized for this resource")
	}

	return next(ctx)
}
