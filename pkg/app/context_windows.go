// +build windows

package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/logger"
)

// Context returns main application context with set logger
// that is canceled when SIGTERM or SIGINT is received.
func Context() context.Context {
	l := zap.L().Named("platform.app")
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.GetCtxWithLogger(ctx, l)

	_ = cancel
	// TODO
	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, unix.SIGTERM, unix.SIGINT)
	// go func() {
	// 	s := <-signals
	// 	signal.Stop(signals)
	// 	l.Sugar().Warnf("Got %s, shutting down...", unix.SignalName(s.(unix.Signal)))
	// 	cancel()
	// }()

	return ctx
}
