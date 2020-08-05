package logger

import (
	"go.uber.org/zap"
)

// SetupGlobalOpts contains logger options.
type SetupGlobalOpts struct {
	LogDebug   bool // enable debug level logging
	LogDevMode bool // enable development mode logging: text instead of JSON, DPanic panics instead of logging errors
}

// SetupGlobal setups global zap logger.
func SetupGlobal(opts *SetupGlobalOpts) {
	if opts == nil {
		opts = new(SetupGlobalOpts)
	}

	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if opts.LogDebug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	if opts.LogDevMode {
		cfg.Development = true
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
}
