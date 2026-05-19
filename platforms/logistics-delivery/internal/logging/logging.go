package logging

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(env string) *Logger {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{logger}
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	fields := []zap.Field{zap.String("service", "logistics-delivery")}
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		fields = append(fields,
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return l.With(fields...)
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Error(msg, fields...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

func (l *Logger) Sync() {
	l.Logger.Sync()
}

type NoopLogger struct{}

func (n *NoopLogger) InfoCtx(_ context.Context, _ string, _ ...zap.Field)  {}
func (n *NoopLogger) ErrorCtx(_ context.Context, _ string, _ ...zap.Field) {}
func (n *NoopLogger) WarnCtx(_ context.Context, _ string, _ ...zap.Field)  {}
func (n *NoopLogger) Sync()                                                 {}
