package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"giautm.dev/awesome/pkg/logging"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"go.uber.org/zap"
)

type authContext struct{}

var (
	authContextKey authContext
)

// TokenFromContext returns the token from the context
func TokenFromContext(ctx context.Context) *jwt.Token {
	token, _ := ctx.Value(authContextKey).(*jwt.Token)
	return token
}

// WithToken returns a new context with the token
func WithToken(ctx context.Context, token *jwt.Token) context.Context {
	if token == nil {
		return ctx
	}

	return context.WithValue(ctx, authContextKey, token)
}

// NewKeyFuncFromEnv returns a keyfunc from the environment
func NewKeyFuncFromEnv() (jwt.Keyfunc, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return nil, errors.New("missing JWT_SECRET in environment variables")
	}

	return KeyFuncHMAC([]byte(key)), nil
}

// KeyFuncHMAC returns a keyfunc for HMAC
func KeyFuncHMAC(secret []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return secret, nil
	}
}

// NewMiddleware returns a new middleware
func NewMiddleware(keyFunc jwt.Keyfunc) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			token, err := extractAndParse(r, keyFunc)
			if err == nil {
				h.ServeHTTP(rw, r.WithContext(WithToken(r.Context(), token)))
				return
			}

			logger := logging.FromContext(r.Context())
			logger.Info("error", zap.Error(err))

			// Should handle Internal Error
			h.ServeHTTP(rw, r)
		})
	}
}

func extractAndParse(r *http.Request, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	rawToken, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
	if err != nil {
		return nil, err
	}

	return jwt.ParseWithClaims(rawToken, &LegacyClaims{}, keyFunc)
}
