package servers

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/tracing"

	"github.com/percona-platform/platform/pkg/logger"
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

	l := zap.L().Named("servers.http").Sugar()

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
			return getCtxForRequest(connCtx)
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

// RequestLoggerMiddleware creates middleware for logging HTTP request execution time.
// It extracts request ID (tracing ID) from incoming HTTP request, creates logger instance with this request ID
// and add logger instance into the request context.
func RequestLoggerMiddleware(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rl := l

		if reqID := tracing.GetRequestIDFromHTTPRequest(r); len(reqID) != 0 {
			rl = l.With(zap.String("request-id", reqID))
		}

		rl = rl.With(zap.String("method", r.Method)).
			With(zap.String("url", r.RequestURI))
		rl.Info("Received request")

		// wrap logger into context so that the following http Handlers could re-use it.
		r = r.WithContext(logger.GetContextWithLogger(r.Context(), rl))
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		rl.Info("Request was processed",
			zap.Int("code", lrw.StatusCode),
			zap.Duration("duration", time.Since(startTime)),
		)
	})
}

// LoggingResponseWriter wrapper struct to catch HTTP response code.
type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// NewLoggingResponseWriter creates wrapper that catches HTTP response code.
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader sends an HTTP response header with the provided
// status code.
// http.ResponseWriter interface implementation.
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
