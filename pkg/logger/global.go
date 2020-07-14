package logger

import (
	"os"

	"go.uber.org/zap"
)

var debug = os.Getenv("PLATFORM_DEBUG") == "1"

// SetupGlobal setups global zap logger.
func SetupGlobal() {
	var err error
	var l *zap.Logger

	if debug {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
}
