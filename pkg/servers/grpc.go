package servers

import (
	"context"
	"crypto/tls"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
)

type GRPCServer interface {
	// Run runs the server until ctx is canceled.
	Run(ctx context.Context)

	// GetUnderlyingServer returns underlying grpc.Server, use it for your server
	// implementation registration. Don't use any control method of returned grpc.Server;
	// use GRPCServer.Run method only.
	GetUnderlyingServer() *grpc.Server
}

type grpcServer struct {
	// TODO remove once we can serve static page with Ingress Controller
	http *http.Server

	grpc *grpc.Server

	addr            string
	shutdownTimeout time.Duration
	l               *zap.SugaredLogger
}

func (s *grpcServer) GetUnderlyingServer() *grpc.Server {
	return s.grpc
}

// stop stops the server.
// It tries to stop the server gracefully until shutdownTimeout is passed, then forcefully.
// Server is fully stopped once method exits.
func (s *grpcServer) stop() {
	// Shutdown returns once HTTP server is stopped, or ctx is canceled (and HTTP server is still running).
	// After that, we forcefully close it anyway.
	httpCtx, httpCancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	go func() {
		if err := s.http.Shutdown(httpCtx); err != nil {
			s.l.Warnf("HTTP Shutdown: %v", err)
		}
		if err := s.http.Close(); err != nil {
			s.l.Warnf("HTTP Close: %v", err)
		}
		httpCancel()
	}()
	<-httpCtx.Done()

	// Since we never call s.grpc.Serve method, gRPC server should be fully stopped by now.
	// Call Stop just to be sure.
	s.grpc.Stop()
}

type NewGRPCServerOpts struct {
	Addr string

	TLSConfig *tls.Config

	// TODO remove once we can serve static page with Ingress Controller
	Handler http.Handler

	WarnDuration    time.Duration
	ShutdownTimeout time.Duration
}

func NewGRPCServer(opts *NewGRPCServerOpts) (GRPCServer, error) {
	l := zap.L().With(zap.String("component", "grpc")).Sugar()

	grpc.EnableTracing = true

	if opts == nil {
		opts = new(NewGRPCServerOpts)
	}

	if opts.Addr == "" {
		l.Panic("No Addr set.")
	}
	if opts.TLSConfig == nil {
		l.Panic("No TLSConfig set.")
	}
	if opts.Handler == nil {
		l.Panic("No Handler set.")
	}

	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 3 * time.Second
	}

	serverOpts := []grpc.ServerOption{
		grpc.ConnectionTimeout(5 * time.Second),
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), //nolint:gomnd

		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryLoggingInterceptor(opts.WarnDuration),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(),
		)),

		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamLoggingInterceptor(opts.WarnDuration),
			grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(),
		)),

		grpc.Creds(credentials.NewTLS(opts.TLSConfig)),
	}
	grpcSrv := grpc.NewServer(serverOpts...)

	httpSrv := &http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if r.ProtoMajor == 2 && mediaType == "application/grpc" {
				grpcSrv.ServeHTTP(rw, r)
				return
			}

			opts.Handler.ServeHTTP(rw, r)
		}),

		TLSConfig: opts.TLSConfig,

		// TODO remove once we have Ingress Controller
		// for now, we need some small values to prevent low and slow attacks
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		ErrorLog: log.New(os.Stderr, "grpc/http.Server", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix),
	}

	return &grpcServer{
		grpc:            grpcSrv,
		http:            httpSrv,
		addr:            opts.Addr,
		shutdownTimeout: opts.ShutdownTimeout,
		l:               l,
	}, nil
}

// Run runs the server until ctx is canceled.
func (s *grpcServer) Run(ctx context.Context) {
	// reflection should not be enabled because we don't want to expose our private APIs
	// reflection.Register(opts.Server)

	channelz.RegisterChannelzServiceToServer(s.grpc)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s.grpc)

	s.l.Infof("Starting server on https://%s/ ...", s.addr)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.l.Panic(err)
	}

	stopped := make(chan struct{})
	go func() {
		<-ctx.Done()
		s.stop()
		close(stopped)
	}()

	err = s.http.ServeTLS(listener, "", "")
	if err == http.ErrServerClosed {
		s.l.Info("Server stopped.")
	} else {
		s.l.Warnf("Server stopped: %v.", err)
	}

	<-stopped

	s.l.Warnf("Listener Close: %v.", listener.Close())
}
