package servers

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
)

// GRPCServer describes common gRPC server interface for Percona Platform applications.
type GRPCServer interface {
	// Run runs the server until ctx is canceled.
	Run(ctx context.Context)

	// GetUnderlyingServer returns underlying grpc.Server, use it for your server
	// implementation registration. Don't use any control method of returned grpc.Server;
	// use GRPCServer.Run method only.
	GetUnderlyingServer() *grpc.Server
}

type grpcServer struct {
	grpc *grpc.Server

	addr            string
	shutdownTimeout time.Duration
	l               *zap.SugaredLogger
}

func (s *grpcServer) GetUnderlyingServer() *grpc.Server {
	return s.grpc
}

// NewGRPCServerOpts configure gRPC server.
type NewGRPCServerOpts struct {
	Addr            string
	WarnDuration    time.Duration
	ShutdownTimeout time.Duration

	// Additional unary and stream interceptors for gRPC server. They will be added at the end of
	// interceptors chain in same order as in this slices. Optional, can be empty.
	ExtraUnaryInterceptors  []grpc.UnaryServerInterceptor
	ExtraStreamInterceptors []grpc.StreamServerInterceptor

	NoAuthMethods []string
}

// NewGRPCServer creates new gRPC server with given options.
func NewGRPCServer(opts *NewGRPCServerOpts) GRPCServer {
	l := zap.L().Named("platform.servers.grpc").Sugar()

	grpc.EnableTracing = true

	if opts == nil {
		opts = new(NewGRPCServerOpts)
	}

	if opts.Addr == "" {
		l.Panic("No Addr set.")
	}

	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 3 * time.Second
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		unaryLoggingInterceptor(opts.WarnDuration),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_validator.UnaryServerInterceptor(),
		unaryAuthInterceptor(opts.NoAuthMethods),
	}
	unaryInterceptors = append(unaryInterceptors, opts.ExtraUnaryInterceptors...)

	streamInterceptors := []grpc.StreamServerInterceptor{
		streamLoggingInterceptor(opts.WarnDuration),
		grpc_prometheus.StreamServerInterceptor,
		grpc_validator.StreamServerInterceptor(),
		streamAuthInterceptor(opts.NoAuthMethods),
	}
	streamInterceptors = append(streamInterceptors, opts.ExtraStreamInterceptors...)

	serverOpts := []grpc.ServerOption{
		grpc.ConnectionTimeout(5 * time.Second),
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), //nolint:gomnd
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	}

	return &grpcServer{
		grpc:            grpc.NewServer(serverOpts...),
		addr:            opts.Addr,
		shutdownTimeout: opts.ShutdownTimeout,
		l:               l,
	}
}

// Run runs the server until ctx is canceled.
// All errors cause panic.
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
		defer close(stopped)
		err := s.grpc.Serve(listener)
		s.l.Infof("Server stopped: %v.", err)
	}()

	<-ctx.Done()

	// try to stop server gracefully, then not
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	go func() {
		<-shutdownCtx.Done()
		s.grpc.Stop()
	}()
	s.grpc.GracefulStop()
	shutdownCancel()

	// listener is already closed there - Serve always closes it on exit,
	// and we can be there only if Serve already exited.
	// But we close it anyway in case gRPC breaks this contract.
	s.l.Infof("Listener closed: %v.", listener.Close())

	<-stopped
}
