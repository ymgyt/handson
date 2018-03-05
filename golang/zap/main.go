package main

import (
	"flag"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOption func(*zap.Config) error

func WithLoggingLevel(level int) LoggerOption {
	return func(cfg *zap.Config) error {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(int8(level)))
		return nil
	}
}

func WithEncoded(encode string) LoggerOption {
	return func(cfg *zap.Config) error {
		switch encode {
		case "json", "console":
		default:
			return fmt.Errorf("unexpected encode %s", encode)
		}
		cfg.Encoding = encode
		return nil
	}
}

func GetLogger(options ...LoggerOption) (*zap.Logger, error) {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(int8(-1))),
		Development:      true,
		Encoding:         "console", // or json
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // CapitalLevelEncoder
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	for _, opt := range options {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	zapOption := zap.AddStacktrace(zapcore.ErrorLevel)
	logger, err := cfg.Build(zapOption)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func run(logger *zap.Logger) {
	logger = logger.With(
		zap.Uint64("request_id", 10),
		zap.String("context", "add context"))

	logger.Debug("debug", zap.String("key", "value"))
	logger.Info("debug", zap.String("key", "value"))
}

func main() {
	var (
		level  int
		encode string
	)
	flag.IntVar(&level, "level", -1, "logging level")
	flag.StringVar(&encode, "encode", "console", "logging encode")
	flag.Parse()

	logger, err := GetLogger(
		WithLoggingLevel(level),
		WithEncoded(encode),
	)
	if err != nil {
		panic(err)
	}

	run(logger)
}
