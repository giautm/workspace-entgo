package server

import (
	"context"
	"net/http"
	"os"

	"giautm.dev/awesome/pkg/logging"
	"giautm.dev/awesome/pkg/ocsentry"
	"github.com/vearutop/sentry-go-exporter-opencensus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gocloud.dev/server"
	"gocloud.dev/server/health"
	"gocloud.dev/server/health/sqlhealth"
)

// NewMuxFx constructs an HTTP mux. Like NewHandler, it depends on *log.Logger.
// However, it also depends on the Fx-specific Lifecycle interface.
//
// A Lifecycle is available in every Fx application. It lets objects hook into
// the application's start and stop phases. In a non-Fx application, the main
// function often includes blocks like this:
//
//   srv, err := NewServer() // some long-running network server
//   if err != nil {
//     log.Fatalf("failed to construct server: %v", err)
//   }
//   // Construct other objects as necessary.
//   go srv.Start()
//   defer srv.Stop()
//
// In this example, the programmer explicitly constructs a bunch of objects,
// crashing the program if any of the constructors encounter unrecoverable
// errors. Once all the objects are constructed, we start any background
// goroutines and defer cleanup functions.
//
// Fx removes the manual object construction with dependency injection. It
// replaces the inline goroutine spawning and deferred cleanups with the
// Lifecycle type.
//
// Here, NewMuxFx makes an HTTP mux available to other functions. Since
// constructors are called lazily, we know that NewMuxFx won't be called unless
// some other function wants to register a handler. This makes it easy to use
// Fx's Lifecycle to start an HTTP server only if we have handlers registered.
func NewMuxFx(lc fx.Lifecycle, logger *zap.Logger, sentryHandler *ocsentry.Handler, sqlCheck *sqlhealth.Checker) *http.ServeMux {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// First, we construct the mux and server. We don't want to start the server
	// until all handlers are registered.
	mux := http.NewServeMux()

	// Create a logger, and assign it to the RequestLogger field of a
	// server.Options struct.
	srvOptions := &server.Options{
		Driver:        server.NewDefaultDriver(),
		HealthChecks:  []health.Checker{sqlCheck},
		RequestLogger: NewRequestLogger(logger),
		TraceExporter: sentry.NewExporter(),
	}

	// Use the constructor function to create the server.
	srv := server.New(sentryHandler.Handle(logging.NewMiddleware(logger)(mux)), srvOptions)

	// If NewMux is called, we know that another function is using the mux. In
	// that case, we'll use the Lifecycle type to register a Hook that starts
	// and stops our HTTP server.
	//
	// Hooks are executed in dependency order. At startup, NewLogger's hooks run
	// before NewMux's. On shutdown, the order is reversed.
	//
	// Returning an error from OnStart hooks interrupts application startup. Fx
	// immediately runs the OnStop portions of any successfully-executed OnStart
	// hooks (so that types which started cleanly can also shut down cleanly),
	// then exits.
	//
	// Returning an error from OnStop hooks logs a warning, but Fx continues to
	// run the remaining hooks.
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 15 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Info("Starting HTTP server.", zap.String("port", port))
			go srv.ListenAndServe(":" + port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server.")

			// Cleanup SQL health check
			sqlCheck.Stop()

			return srv.Shutdown(ctx)
		},
	})

	return mux
}
