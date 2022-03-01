package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

// Guard is the user guard
type Guard string

const (
	// GuardEmployer is the guard for employer
	GuardEmployer Guard = "employer"
	// GuardEmployee is the guard for worker
	GuardWorker Guard = "worker"
)

// LegacyClaims is the legacy claims
type LegacyClaims struct {
	jwt.StandardClaims

	ID        int    `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Guard     Guard  `json:"guard"`
}

// ClaimsFromContext returns the claims from the context
func ClaimsFromContext(ctx context.Context) *LegacyClaims {
	token := TokenFromContext(ctx)
	if token == nil {
		return nil
	}

	if c, ok := token.Claims.(*LegacyClaims); ok {
		return c
	}

	return nil
}
