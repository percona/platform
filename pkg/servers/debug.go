package servers

import (
	"bytes"
	"context"
	_ "expvar" // register /debug/vars
	"log"
	"net/http"
	_ "net/http/pprof" // register /debug/pprof
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Percona-Platform/platform/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type RunDebugServerOpts struct {
	Addr            string
	ShutdownTimeout time.Duration
}

func RunDebugServer(ctx context.Context, opts *RunDebugServerOpts) {
	if opts == nil {
		opts = new(RunDebugServerOpts)
	}

	l := zap.L().With(zap.String("component", "debug")).Sugar()

	if opts.Addr == "" {
		l.Panic("No Addr set,")
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 3 * time.Second
	}

	handler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		ErrorLog:      &logger.PromHTTP{L: l},
		ErrorHandling: promhttp.ContinueOnError,
	})
	http.Handle("/debug/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handler))

	handlers := []string{
		"/debug/metrics",  // by http.Handle above
		"/debug/vars",     // by expvar
		"/debug/requests", // by golang.org/x/net/trace imported by google.golang.org/grpc
		"/debug/events",   // by golang.org/x/net/trace imported by google.golang.org/grpc
		"/debug/pprof",    // by net/http/pprof
	}
	for i, h := range handlers {
		handlers[i] = "http://" + opts.Addr + h
	}

	var buf bytes.Buffer
	err := template.Must(template.New("debug").Parse(`
	<html>
	<body>
	<ul>
	{{ range . }}
		<li><a href="{{ . }}">{{ . }}</a></li>
	{{ end }}
	</ul>
	</body>
	</html>
	`)).Execute(&buf, handlers)
	if err != nil {
		l.Panic(err)
	}
	http.HandleFunc("/debug", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(buf.Bytes())
	})
	l.Infof("Starting server on http://%s/debug\nRegistered handlers:\n\t%s", opts.Addr, strings.Join(handlers, "\n\t"))

	server := &http.Server{
		Addr:     opts.Addr,
		ErrorLog: log.New(os.Stderr, "runDebugServer: ", 0),
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			l.Panic(err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), opts.ShutdownTimeout)
	if err := server.Shutdown(ctx); err != nil {
		l.Errorf("Failed to shutdown gracefully: %s", err)
	}
	cancel()
}
