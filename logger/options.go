package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigOption interface {
	apply(cfg *zap.Config)
}

type optionFunc func(cfg *zap.Config)

func (f optionFunc) apply(cfg *zap.Config) {
	f(cfg)
}

func FormatTime(layout string) ConfigOption {
	return optionFunc(func(k *zap.Config) {
		k.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(layout)
	})
}

func LogLevel(level Level) ConfigOption {
	return optionFunc(func(cfg *zap.Config) {
		cfg.Level.SetLevel(level)
	})
}
