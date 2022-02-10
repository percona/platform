//go:build !windows
// +build !windows

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
	// catch the common service initialization problem
	if !logger.FlagsParsed {
		panic("app.Context should be called after app.Setup and kingpin.Parse")
	}

	l := zap.L()
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.GetContextWithLogger(ctx, l)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGTERM, unix.SIGINT)
	go func() {
		s := <-signals
		signal.Stop(signals)
		l.Sugar().Warnf("Got %s, shutting down...", unix.SignalName(s.(unix.Signal))) // nolint: forcetypeassert
		cancel()
	}()

	return ctx
}
