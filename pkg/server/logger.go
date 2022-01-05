package server

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gocloud.dev/server/requestlog"
)

type zapRequestLogger struct {
	logger *zap.Logger
}

func (rl *zapRequestLogger) Log(ent *requestlog.Entry) {
	var remoteIP zapcore.Field
	if ent.RemoteIP != "" {
		remoteIP = zap.String("remoteIP", ent.RemoteIP)
	} else {
		remoteIP = zap.String("remoteIP", "-")
	}

	rl.logger.Debug("request", remoteIP,
		zap.Time("receivedTime", ent.ReceivedTime),
		zap.String("requestMethod", ent.RequestMethod),
		zap.String("requestURL", ent.RequestURL),
		zap.String("proto", ent.Proto),
		zap.Int("status", ent.Status),
		zap.Int64("responseBodySize", ent.ResponseBodySize),
		zap.String("referer", ent.Referer),
		zap.String("userAgent", ent.UserAgent))
}
