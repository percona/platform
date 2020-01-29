package logger

import (
	"go.uber.org/zap"
)

func SetupGlobal() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
}
