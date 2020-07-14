package logger

import (
	"os"

	"go.uber.org/zap"
)

// SetupGlobal setups global zap logger.
func SetupGlobal() {
	var err error
	var l *zap.Logger

	if os.Getenv("PLATFORM_DEBUG") == "1" {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
}
