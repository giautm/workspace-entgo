package sentry

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func NewSentry() (*sentryhttp.Handler, error) {
	err := sentry.Init(sentry.ClientOptions{
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return nil, fmt.Errorf("sentry: failed to init: %w", err)
	}

	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}), nil
}
