// +build !windows

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/logger"
)

// Context returns main application context with set logger
// that is canceled when SIGTERM or SIGINT is received.
func Context() context.Context {
	l := zap.L().Named("platform.app")
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.GetCtxWithLogger(ctx, l)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-signals
		signal.Stop(signals)
		l.Sugar().Warnf("Got %s, shutting down...", s.(syscall.Signal).String())
		cancel()
	}()

	return ctx
}
