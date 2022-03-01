package logger

import (
	"context"
	"os"

	"giautm.dev/awesome/internal/buildinfo"
	"giautm.dev/awesome/pkg/logging"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// Module exports the logger module.
var Module = fx.Options(
	fx.Provide(NewLoggerFx),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)

// NewLoggerFx returns a new logger.
func NewLoggerFx(lc fx.Lifecycle) (*zap.Logger, error) {
	logger, err := logging.NewLoggerFromEnv()
	if err != nil {
		return nil, err
	}
	logger = logger.With(
		zap.String("buildID", buildinfo.BuildID),
		zap.String("buildTag", buildinfo.BuildTag),
		zap.String("podName", os.Getenv("K8S_POD_NAME")),
	)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			// NOTE(giautm): Ignore error from log.Sync()
			// See: https://github.com/uber-go/zap/issues/880
			logger.Sync()
			return nil
		},
	})

	return logger, nil
}
