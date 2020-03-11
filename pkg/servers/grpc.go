package servers

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"

	"github.com/percona-platform/platform/pkg/ptls"
)

type GetGRPCServerOpts struct {
	Cert string
	Key  string

	CertFile string
	KeyFile  string

	ACME *ptls.GetACMEOpts

	Handler      http.Handler
	WarnDuration time.Duration
}

type GRPCServer struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

func (s *GRPCServer) Serve(listener net.Listener) error {
	if s.HTTP != nil {
		return s.HTTP.ServeTLS(listener, "", "")
	}

	return s.GRPC.Serve(listener)
}

func (s *GRPCServer) Stop(listener net.Listener) {
	if s.HTTP != nil {
		s.HTTP.Close()
	}

	s.GRPC.Stop()
}

func (s *GRPCServer) GracefulStop(timeout time.Duration) {
	// try to stop server gracefully, then not
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if s.HTTP != nil {
		s.HTTP.Shutdown(ctx)
	}

	go func() {
		<-ctx.Done()
		s.GRPC.Stop()
	}()
	s.GRPC.GracefulStop()
}

func GetGRPCServer(opts *GetGRPCServerOpts) (*GRPCServer, http.Handler, error) {
	l := zap.L().With(zap.String("component", "grpc")).Sugar()

	grpc.EnableTracing = true

	if opts == nil {
		opts = new(GetGRPCServerOpts)
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
	}

	var tlsConfig *tls.Config
	var handler http.Handler
	var err error
	switch {
	case opts.Cert != "" && opts.Key != "":
		if opts.CertFile != "" || opts.KeyFile != "" {
			return nil, nil, errors.New("both Cert/Key and CertFile/KeyFile are specified")
		}
		if opts.ACME != nil {
			return nil, nil, errors.New("both Cert/Key and ACME are specified")
		}

		l.Info("Using given certificate and key for gRPC server.")

		cert, err := tls.X509KeyPair([]byte(opts.Cert), []byte(opts.Key))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to parse TLS data")
		}

		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	case opts.CertFile != "" && opts.KeyFile != "":
		if opts.ACME != nil {
			return nil, nil, errors.New("both CertFile/KeyFile and ACME are specified")
		}

		l.Info("Using given certificate and key files for gRPC server.")

		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to load TLS files")
		}

		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	case opts.ACME != nil:
		l.Infof("Using ACME (%v) for gRPC server.", opts.ACME.Hosts)

		tlsConfig, handler, err = ptls.GetACME(opts.ACME)
		if err != nil {
			return nil, nil, err
		}
	}

	if tlsConfig != nil {
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	grpcServer := grpc.NewServer(serverOpts...)

	var httpServer *http.Server
	if opts.Handler != nil {
		httpServer = &http.Server{
			Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(rw, r)
					return
				}

				opts.Handler.ServeHTTP(rw, r)
			}),
		}

		if tlsConfig != nil {
			httpServer.TLSConfig = tlsConfig
		}
	}

	return &GRPCServer{GRPC: grpcServer, HTTP: httpServer}, handler, nil
}

type RunGRPCServerOpts struct {
	Server          *GRPCServer
	Addr            string
	ShutdownTimeout time.Duration
}

func RunGRPCServer(ctx context.Context, opts *RunGRPCServerOpts) {
	if opts == nil {
		opts = new(RunGRPCServerOpts)
	}

	l := zap.L().With(zap.String("component", "grpc")).Sugar()

	if opts.Server == nil {
		l.Panic("No Server set.")
	}
	if opts.Addr == "" {
		l.Panic("No Addr set.")
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 3 * time.Second
	}

	// reflection should not be enabled because we don't want to expose our private APIs
	// reflection.Register(opts.Server)

	channelz.RegisterChannelzServiceToServer(opts.Server.GRPC)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(opts.Server.GRPC)

	// run server until it is stopped gracefully or not
	l.Infof("Starting server on https://%s/ ...", opts.Addr)
	listener, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		l.Info("Server started.")
		var err error
		for {
			err = opts.Server.Serve(listener)
			if err == nil || err == grpc.ErrServerStopped || err == http.ErrServerClosed {
				break
			}
			l.Errorf("Failed to serve: %s", err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()

	opts.Server.GracefulStop(opts.ShutdownTimeout)
}
