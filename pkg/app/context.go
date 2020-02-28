package app

import (
	"context"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"golang.org/x/sys/unix"

	"github.com/percona-platform/platform/pkg/logger"
)

// Context returns main application context with set logger
// that is canceled when SIGTERM or SIGINT is received.
func Context() context.Context {
	l := zap.L().With(zap.String("component", "app"))
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.GetCtxWithLogger(ctx, l)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGTERM, unix.SIGINT)
	go func() {
		s := <-signals
		signal.Stop(signals)
		l.Sugar().Warnf("Got %s, shutting down...", unix.SignalName(s.(unix.Signal)))
		cancel()
	}()

	return ctx
}
