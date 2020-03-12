package logger

import (
	"io"

	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

// GRPC is a compatibility wrapper between zap's sugared logger entry and gRPC logger interface.
type GRPC struct {
	*zap.SugaredLogger

	// Set to true for very verbose gRPC logging.
	Verbose bool
}

// V reports whether verbosity level l is at least the requested verbose level.
func (g *GRPC) V(l int) bool {
	return g.Verbose
}

func (g *GRPC) Warningf(format string, args ...interface{}) {
	// Inhibit the very specific message that spams our logs due to AWS NLB configuration for bare-bones telemetry:
	// https://github.com/grpc/grpc-go/blob/142182889d38b76209f1d9f1d8e91d7608aff542/server.go#L685
	// Remove this hack once we have "normal" load balancer / ingress controller.
	if format == "grpc: Server.Serve failed to complete security handshake from %q: %v" && len(args) == 2 && args[1] == io.EOF {
		return
	}

	g.Warnf(format, args...)
}

func (g *GRPC) Infoln(args ...interface{})    { g.Info(args...) }
func (g *GRPC) Warning(args ...interface{})   { g.Warn(args...) }
func (g *GRPC) Warningln(args ...interface{}) { g.Warn(args...) }
func (g *GRPC) Errorln(args ...interface{})   { g.Error(args...) }
func (g *GRPC) Fatalln(args ...interface{})   { g.Fatal(args...) }

// check interfaces
var (
	_ grpclog.LoggerV2 = (*GRPC)(nil)
)
