package sentry

import (
	"context"
	"fmt"
	"time"

	"giautm.dev/awesome/internal/project"
	"giautm.dev/awesome/pkg/ocsentry"
	"github.com/getsentry/sentry-go"
	"go.opencensus.io/trace"
	"go.uber.org/fx"
)

// NewSentryFx is provider for sentry handler
func NewSentryFx(lc fx.Lifecycle) (*ocsentry.Handler, error) {
	clientOpts := sentry.ClientOptions{
		TracesSampleRate: 0.2,
	}

	if project.DevMode() {
		clientOpts.TracesSampleRate = 1.0
		trace.ApplyConfig(trace.Config{
			DefaultSampler: trace.ProbabilitySampler(1.0),
		})
	}

	if err := sentry.Init(clientOpts); err != nil {
		return nil, fmt.Errorf("sentry: failed to init: %w", err)
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			sentry.Flush(time.Second)
			return nil
		},
	})

	httpOpts := ocsentry.Options{
		Repanic: false,
	}
	if project.DevMode() {
		httpOpts.Repanic = true
	}

	return ocsentry.New(httpOpts), nil
}
