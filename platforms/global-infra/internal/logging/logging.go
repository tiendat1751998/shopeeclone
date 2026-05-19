package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(appName, logLevel string) *zap.Logger {
	level := zapcore.InfoLevel
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(level)
	cfg.InitialFields = map[string]interface{}{"app": appName}
	logger, _ := cfg.Build()
	return logger
}
