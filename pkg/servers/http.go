package servers

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// RunHTTPServerOpts configure HTTP server.
type RunHTTPServerOpts struct {
	Addr            string
	Handler         http.Handler
	ShutdownTimeout time.Duration
}

// RunHTTPServer runs HTTP server with given options until ctx is canceled.
// All errors cause panic.
func RunHTTPServer(ctx context.Context, opts *RunHTTPServerOpts) {
	if opts == nil {
		opts = new(RunHTTPServerOpts)
	}

	l := zap.L().Named("platform.servers.http").Sugar()

	if opts.Addr == "" {
		l.Panic("No Addr set.")
	}
	if opts.Handler == nil {
		opts.Handler = http.NotFoundHandler()
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 3 * time.Second
	}

	l.Infof("Starting server on http://%s/", opts.Addr)

	server := &http.Server{
		Addr: opts.Addr,
		ErrorLog: log.New(
			os.Stderr,
			"platform.servers.http.Server",
			log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix,
		),

		// propagate ctx cancellation signals to handlers
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},

		// propagate ctx cancellation signals and pass logger to handlers
		ConnContext: func(connCtx context.Context, _ net.Conn) context.Context {
			c, _ := getCtxForRequest(connCtx)
			return c
		},

		Handler: opts.Handler,
	}

	stopped := make(chan error)
	go func() {
		stopped <- server.ListenAndServe()
	}()

	// any ListenAndServe error before ctx is canceled is fatal
	select {
	case <-ctx.Done():
	case err := <-stopped:
		l.Panicf("Unexpected server stop: %v.", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), opts.ShutdownTimeout)
	if err := server.Shutdown(shutdownCtx); err != nil {
		l.Errorf("Failed to shutdown gracefully: %s", err)
	}
	shutdownCancel()

	<-stopped
	l.Info("Server stopped.")
}
