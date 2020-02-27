package servers

import (
	"context"
	"runtime/debug"
	"runtime/pprof"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	validator "github.com/mwitkow/go-proto-validators"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona-platform/platform/pkg/logger"
)

func logRequest(l *zap.Logger, prefix string, warnD time.Duration, f func() error) (err error) {
	start := time.Now()
	sl := l.Sugar()
	sl.Infof("Starting %s ...", prefix)

	defer func() {
		dur := time.Since(start)

		if p := recover(); p != nil {
			// Always log with %+v, even before re-panic - there can be inner stacktraces
			// produced by panic(errors.WithStack(err)).
			// Also always log debug.Stack() for all panics.
			sl.DPanicf("%s done in %s with panic: %+v\nStack: %s", prefix, dur, p, debug.Stack())

			err = status.Error(codes.Internal, "Internal server error.")
			return
		}

		// log gRPC errors as warning, not errors, even if they are wrapped
		_, gRPCError := status.FromError(errors.Cause(err))
		switch {
		case err == nil:
			if warnD == 0 || dur < warnD {
				sl.Infof("%s done in %s.", prefix, dur)
			} else {
				sl.Warnf("%s done in %s (quite long).", prefix, dur)
			}
		case gRPCError:
			// %+v for inner stacktraces produced by errors.WithStack(err)
			sl.Warnf("%s done in %s with gRPC error: %+v", prefix, dur, err)
		default:
			// %+v for inner stacktraces produced by errors.WithStack(err)
			sl.Errorf("%s done in %s with unexpected error: %+v", prefix, dur, err)
			err = status.Error(codes.Internal, "Internal server error.")
		}
	}()

	err = f()
	return //nolint:nakedret
}

func validate(req interface{}) error {
	if v, ok := req.(validator.Validator); ok {
		if err := v.Validate(); err != nil {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	return nil
}

type unary struct {
	warnDuration time.Duration
}

// intercept adds pprof labels, context logger, validation, and Prometheus metrics to unary server RPC.
func (u unary) intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// add pprof labels for more useful profiles
	defer pprof.SetGoroutineLabels(ctx)
	ctx = pprof.WithLabels(ctx, pprof.Labels("method", info.FullMethod))
	pprof.SetGoroutineLabels(ctx)

	// set logger
	l := zap.L().With(zap.String("request", logger.MakeRequestID()))
	ctx = logger.SetEntry(ctx, l)

	var res interface{}
	err := logRequest(l, "RPC "+info.FullMethod, u.warnDuration, func() error {
		var origErr error

		if origErr = validate(req); origErr != nil {
			return origErr
		}

		res, origErr = grpc_prometheus.UnaryServerInterceptor(ctx, req, info, handler)
		return origErr
	})

	// err is already logged by logRequest
	l.Sugar().Debugf("\nRequest:\n%s\nResponse:\n%s\n", req, res)

	return res, err
}

// vServerStream is a thin wrapper around grpc.ServerStream that wraps RecvMsg with validation.
// TODO Wrap SendMsg, RecvMsg with logging once we need streams (?)
type vServerStream struct {
	grpc.ServerStream
}

func (s *vServerStream) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	return validate(m)
}

type stream struct {
	warnDuration time.Duration
}

// intercept adds pprof labels, context logger, validation, and Prometheus metrics to stream server RPC.
func (s stream) intercept(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	// add pprof labels for more useful profiles
	defer pprof.SetGoroutineLabels(ctx)
	ctx = pprof.WithLabels(ctx, pprof.Labels("method", info.FullMethod))
	pprof.SetGoroutineLabels(ctx)

	// set logger
	l := zap.L().With(zap.String("request", logger.MakeRequestID()))
	ctx = logger.SetEntry(ctx, l)

	err := logRequest(l, "Stream "+info.FullMethod, s.warnDuration, func() error {
		ss = &vServerStream{ss}

		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx

		return grpc_prometheus.StreamServerInterceptor(srv, wrapped, info, handler)
	})
	return err
}

// check interfaces
var (
	_ grpc.UnaryServerInterceptor  = new(unary).intercept
	_ grpc.StreamServerInterceptor = new(stream).intercept
)
