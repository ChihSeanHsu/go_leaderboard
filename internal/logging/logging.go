package logging

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogging() {
	// See the documentation for Config and zapcore.EncoderConfig for all the available options.
	// https://godoc.org/go.uber.org/zap/zapcore#EncoderConfig
	logger, err := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:          "json",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: false,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}.Build()
	if err != nil {
		panic(err)
	}

	// make it a global logger
	zap.ReplaceGlobals(logger)
}

// alias to zap.S() so that caller doesn't need to be aware of zap
func ZapL(ctx ...context.Context) *zap.SugaredLogger {
	k := "TraceID"
	if len(ctx) == 0 {
		return zap.S()
	} else if v := ctx[0].Value(k); v != nil {
		return zap.S().With(k, v)
	} else {
		return zap.S()
	}
}
