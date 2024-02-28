package logging

import (
	"git.stamus-networks.com/lanath/stamus-ctl/internal/app"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	envType = "dev"
	Sugar   *zap.SugaredLogger
	levels  = [...]zapcore.Level{zap.WarnLevel, zap.InfoLevel, zap.DebugLevel}
)

func SetLogger() {
	verbosity := viper.GetInt("verbose")
	encoder := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if verbosity >= len(levels) {
		verbosity = len(levels) - 1
	}

	if envType == "prd" {
		encoder.StacktraceKey = zapcore.OmitKey

		if app.Name == "stamusctl" {
			encoder.TimeKey = zapcore.OmitKey
			encoder.LevelKey = zapcore.OmitKey
			encoder.CallerKey = zapcore.OmitKey
		}
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(levels[verbosity]),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoder,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, _ := config.Build()
	defer log.Sync()
	Sugar = log.Sugar()
}

func init() {
	noop := zap.NewNop()

	Sugar = noop.Sugar()
}
