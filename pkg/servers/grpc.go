package servers

import (
	"context"
	"crypto/tls"
	"mime"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
)

type GetGRPCServerOpts struct {
	Addr string

	Cert string
	Key  string

	CertFile string
	KeyFile  string

	TLSConfig *tls.Config

	Handler         http.Handler
	WarnDuration    time.Duration
	ShutdownTimeout time.Duration
}

type GRPCServer struct {
	grpc            *grpc.Server
	http            *http.Server
	addr            string
	shutdownTimeout time.Duration
}

// RegisterGRPCServer allows to register GRPC server implementation in
// underlying grpc.Server
func (s *GRPCServer) RegisterGRPCServer(f func(s *grpc.Server)) {
	f(s.grpc)
}

func (s *GRPCServer) Serve(listener net.Listener) error {
	if s.http != nil {
		return s.http.ServeTLS(listener, "", "")
	}

	return s.grpc.Serve(listener)
}

func (s *GRPCServer) Stop(listener net.Listener) {
	if s.http != nil {
		s.http.Close()
	}

	s.grpc.Stop()
}

func (s *GRPCServer) GracefulStop(timeout time.Duration) {
	// try to stop server gracefully, then not
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if s.http != nil {
		s.http.Shutdown(ctx)
	}

	go func() {
		<-ctx.Done()
		s.grpc.Stop()
	}()
	s.grpc.GracefulStop()
}

func NewGRPCServer(opts *GetGRPCServerOpts) (*GRPCServer, error) {
	l := zap.L().With(zap.String("component", "grpc")).Sugar()

	grpc.EnableTracing = true

	if opts == nil {
		opts = new(GetGRPCServerOpts)
	}

	if opts.Addr == "" {
		l.Panic("No Addr set.")
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
	}

	var tlsConfig *tls.Config
	switch {
	case opts.Cert != "" && opts.Key != "":
		if opts.CertFile != "" || opts.KeyFile != "" {
			return nil, errors.New("both Cert/Key and CertFile/KeyFile are specified")
		}
		if opts.TLSConfig != nil {
			return nil, errors.New("both Cert/Key and TLSConfig are specified")
		}

		l.Info("Using given certificate and key for gRPC server.")

		cert, err := tls.X509KeyPair([]byte(opts.Cert), []byte(opts.Key))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse TLS data")
		}

		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	case opts.CertFile != "" && opts.KeyFile != "":
		if opts.TLSConfig != nil {
			return nil, errors.New("both CertFile/KeyFile and TLSConfig are specified")
		}

		l.Info("Using given certificate and key files for gRPC server.")

		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load TLS files")
		}

		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	case opts.TLSConfig != nil:
		l.Infof("Using TLSConfig for gRPC server.")
		tlsConfig = opts.TLSConfig
	}

	if tlsConfig != nil {
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	grpcServer := grpc.NewServer(serverOpts...)

	var httpServer *http.Server
	if opts.Handler != nil {
		httpServer = &http.Server{
			Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if r.ProtoMajor == 2 && mediaType == "application/grpc" {
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

	return &GRPCServer{
		grpc: grpcServer,
		http: httpServer,
		addr: opts.Addr,
	}, nil
}

func (s *GRPCServer) Start(ctx context.Context) {
	l := zap.L().With(zap.String("component", "grpc")).Sugar()

	// reflection should not be enabled because we don't want to expose our private APIs
	// reflection.Register(opts.Server)

	channelz.RegisterChannelzServiceToServer(s.grpc)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s.grpc)

	// run server until it is stopped gracefully or not
	l.Infof("Starting server on https://%s/ ...", s.addr)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	go func() {
		l.Info("Server started.")
		for {
			err = s.Serve(listener)
			if err == nil || err == grpc.ErrServerStopped || err == http.ErrServerClosed {
				break
			}
			l.Errorf("Failed to serve: %s", err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()

	s.GracefulStop(s.shutdownTimeout)
}
