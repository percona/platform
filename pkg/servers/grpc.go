package servers

import (
	"context"
	"crypto/tls"
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
	"google.golang.org/grpc/reflection"

	"github.com/percona-platform/platform/pkg/ptls"
)

type GetGRPCServerOpts struct {
	Cert string
	Key  string

	CertFile string
	KeyFile  string

	ACME *ptls.GetACMEOpts

	WarnDuration time.Duration
}

func GetGRPCServer(opts *GetGRPCServerOpts) (*grpc.Server, http.Handler, error) {
	grpc.EnableTracing = true

	if opts == nil {
		opts = new(GetGRPCServerOpts)
	}

	serverOpts := []grpc.ServerOption{
		grpc.ConnectionTimeout(5 * time.Second),
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), //nolint:gomnd

		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			NewUnaryLoggingInterceptor(opts.WarnDuration),
			grpc_validator.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			NewStreamLoggingInterceptor(opts.WarnDuration),
			grpc_validator.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		)),
	}

	var creds credentials.TransportCredentials
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

		cert, err := tls.X509KeyPair([]byte(opts.Cert), []byte(opts.Key))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to parse TLS data")
		}

		creds = credentials.NewServerTLSFromCert(&cert)

	case opts.CertFile != "" && opts.KeyFile != "":
		if opts.ACME != nil {
			return nil, nil, errors.New("both CertFile/KeyFile and ACME are specified")
		}

		creds, err = credentials.NewServerTLSFromFile(opts.CertFile, opts.KeyFile)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to load TLS files")
		}

	case opts.ACME != nil:
		var tlsConfig *tls.Config
		tlsConfig, handler, err = ptls.GetACME(opts.ACME)
		if err != nil {
			return nil, nil, err
		}

		creds = credentials.NewTLS(tlsConfig)
	}

	if creds != nil {
		serverOpts = append(serverOpts, grpc.Creds(creds))
	}

	return grpc.NewServer(serverOpts...), handler, nil
}

type RunGRPCServerOpts struct {
	Server          *grpc.Server
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

	reflection.Register(opts.Server)

	channelz.RegisterChannelzServiceToServer(opts.Server)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(opts.Server)

	// run server until it is stopped gracefully or not
	l.Infof("Starting server on https://%s/ ...", opts.Addr)
	listener, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		l.Info("Server started.")
		for {
			err = opts.Server.Serve(listener)
			if err == nil || err == grpc.ErrServerStopped {
				break
			}
			l.Errorf("Failed to serve: %s", err)
		}
		l.Info("Server stopped.")
	}()

	<-ctx.Done()

	// try to stop server gracefully, then not
	ctx, cancel := context.WithTimeout(context.Background(), opts.ShutdownTimeout)
	go func() {
		<-ctx.Done()
		opts.Server.Stop()
	}()
	opts.Server.GracefulStop()
	cancel()
}
