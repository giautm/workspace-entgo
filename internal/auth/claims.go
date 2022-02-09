package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type Guard string

const (
	GuardEmployer Guard = "employer"
	GuardWorker   Guard = "worker"
)

type LegacyClaims struct {
	jwt.StandardClaims

	ID        int    `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Guard     Guard  `json:"guard"`
}

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
