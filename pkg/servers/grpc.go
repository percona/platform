package servers

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/Percona-Platform/platform/pkg/ptls"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type GetGRPCServerOpts struct {
	CertFile     string
	KeyFile      string
	ACME         *ptls.GetACMEOpts
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
			unary{
				warnDuration: opts.WarnDuration,
			}.intercept,
			grpc_validator.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			stream{
				warnDuration: opts.WarnDuration,
			}.intercept,
			grpc_validator.StreamServerInterceptor(),
		)),
	}

	var creds credentials.TransportCredentials
	var handler http.Handler
	var err error
	switch {
	case opts.CertFile != "" && opts.KeyFile != "":
		if opts.ACME != nil {
			return nil, nil, errors.New("both CertFile/KeyFile and ACME are specified")
		}

		creds, err = credentials.NewServerTLSFromFile("dev-cert.pem", "dev-key.pem")
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
