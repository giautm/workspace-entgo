package main

import (
	"net/http"

	_ "giautm.dev/awesome/ent/runtime"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"giautm.dev/awesome/graphql"
	"giautm.dev/awesome/internal/auth"
	"giautm.dev/awesome/internal/database"
	"giautm.dev/awesome/internal/logger"
	"giautm.dev/awesome/internal/sentry"
	"giautm.dev/awesome/pkg/server"
)

// NewHandler constructs a simple HTTP handler. Since it returns an
// http.Handler, Fx will treat NewHandler as the constructor for the
// http.Handler type.
func NewHandler(logger *zap.Logger, keyFunc jwt.Keyfunc, gqlserver *graphql.Server) (http.Handler, error) {
	logger.Info("Executing NewHandler.")

	r := chi.NewRouter()
	r.Use(auth.NewMiddleware(keyFunc))
	r.Route("/", gqlserver.MountRoutes)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	return c.Handler(r), nil
}

// Register mounts our HTTP handler on the mux.
func Register(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}

func main() {
	app := fx.New(
		logger.Module,
		database.Module,
		graphql.Module,
		fx.Provide(
			sentry.NewSentryFx,
			auth.NewKeyFuncFromEnv,
			NewHandler,
			server.NewMuxFx,
		),
		// Since constructors are called lazily, we need some invocations to
		// kick-start our application. In this case, we'll use Register. Since it
		// depends on an http.Handler and *http.ServeMux, calling it requires Fx
		// to build those types using the constructors above. Since we call
		// NewMux, we also register Lifecycle hooks to start and stop an HTTP
		// server.
		fx.Invoke(Register),
	)
	app.Run()
}
