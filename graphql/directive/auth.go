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

// AuthRole is the role that can be used to restrict access to a field
type AuthRole string

const (
	// AuthRoleAdmin is the role that can access all fields
	AuthRoleAdmin AuthRole = "ADMIN"
)

// AllRole is the list of valid roles
var AllRole = []AuthRole{
	AuthRoleAdmin,
}

// IsValid returns true if the role is valid
func (e AuthRole) IsValid() bool {
	switch e {
	case AuthRoleAdmin:
		return true
	}
	return false
}

// String implements the Stringer interface
func (e AuthRole) String() string {
	return string(e)
}

// UnmarshalRole implements the GQL Unmarshaler interface
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

// MarshalRole implements the GQL Marshaler interface
func (e AuthRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Auth is the directive that can be used to restrict access to a field
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
