package observability

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func InitLogger(serviceName, level string) *zap.Logger {
	once.Do(func() {
		var cfg zap.Config

		if os.Getenv("APP_ENV") == "production" {
			cfg = zap.NewProductionConfig()
			cfg.EncoderConfig.TimeKey = "timestamp"
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		switch level {
		case "debug":
			cfg.Level.SetLevel(zapcore.DebugLevel)
		case "info":
			cfg.Level.SetLevel(zapcore.InfoLevel)
		case "warn":
			cfg.Level.SetLevel(zapcore.WarnLevel)
		case "error":
			cfg.Level.SetLevel(zapcore.ErrorLevel)
		default:
			cfg.Level.SetLevel(zapcore.InfoLevel)
		}

		cfg.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
		cfg.EncoderConfig.CallerKey = "caller"

		l, err := cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}

		logger = l.With(zap.String("service", serviceName))
	})

	return logger
}

func GetLogger() *zap.Logger {
	if logger == nil {
		return InitLogger("unknown", "info")
	}
	return logger
}

func LogWithTrace(ctx context.Context) *zap.Logger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return GetLogger().With(
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return GetLogger()
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}
