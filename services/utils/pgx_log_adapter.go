package log

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PgxLogAdapter struct{}

func (PgxLogAdapter) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	logger := FromContext(ctx).Desugar().
		WithOptions(zap.AddCallerSkip(5)).
		WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	fields := make([]zapcore.Field, 0, len(data))
	for k, v := range data {
		if k == "time" {
			// below renaming is required by GKE to properly display SQL queries
			k = "queryTime"
		}
		fields = append(fields, zap.Reflect(k, v))
	}

	switch level {
	case pgx.LogLevelTrace:
		logger.Debug(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelDebug:
		logger.Debug(msg, fields...)
	case pgx.LogLevelInfo:
		logger.Info(msg, fields...)
	case pgx.LogLevelWarn:
		logger.Warn(msg, fields...)
	case pgx.LogLevelError:
		logger.Error(msg, fields...)
	default:
		logger.Error(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	}
}
