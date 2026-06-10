package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New builds a structured zap logger. In production it emits JSON; otherwise a
// human-friendly console encoder is used. level accepts "debug"|"info"|"warn"|"error".
func New(level string, production bool) (*zap.Logger, error) {
	lvl := zapcore.InfoLevel
	_ = lvl.UnmarshalText([]byte(level))

	var cfg zap.Config
	if production {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build(zap.AddCallerSkip(0))
}
