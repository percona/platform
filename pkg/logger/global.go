package logger

import (
	"go.uber.org/zap"
)

// SetupGlobal setups global zap logger.
func SetupGlobal(logDebug, logDevMode bool) {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if logDebug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	if logDevMode {
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
