package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	var (
		logLevel = flag.Int("logLevel", 0, "zap log level debug(-1) info(0) warn(1)")
	)
	flag.Parse()

	logger, err := GetLogger(*logLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&UpdateCmd{log: logger}, "")

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

func GetLogger(level int) (*zap.Logger, error) {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(int8(level))),
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
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	option := zap.AddStacktrace(zapcore.ErrorLevel)
	return cfg.Build(option)
}
