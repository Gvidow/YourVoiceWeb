package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

type Level = zapcore.Level

func New(options ...ConfigOption) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	for _, option := range options {
		option.apply(&cfg)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("new Logger: %w", err)
	}
	return &Logger{logger}, nil
}
