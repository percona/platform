package servers

import (
	"context"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
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

// intercept adds context logger and Prometheus metrics to unary server RPC.
func NewUnaryLoggingInterceptor(warnDuration time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// add pprof labels for more useful profiles
		defer pprof.SetGoroutineLabels(ctx)
		ctx = pprof.WithLabels(ctx, pprof.Labels("method", info.FullMethod))
		pprof.SetGoroutineLabels(ctx)

		// set logger
		l := zap.L().With(zap.String("request", logger.MakeRequestID()))
		ctx = logger.SetEntry(ctx, l)

		var res interface{}
		err := logRequest(l, "RPC "+info.FullMethod, warnDuration, func() error {
			var origErr error

			res, origErr = handler(ctx, req)
			// l.Debugf("\nRequest:\n%s\nResponse:\n%s\n", req, res)
			return origErr
		})
		return res, err
	}
}

// Stream adds context logger and Prometheus metrics to stream server RPC.
func NewStreamLoggingInterceptor(warnDuration time.Duration) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		// add pprof labels for more useful profiles
		defer pprof.SetGoroutineLabels(ctx)
		ctx = pprof.WithLabels(ctx, pprof.Labels("method", info.FullMethod))
		pprof.SetGoroutineLabels(ctx)

		// set logger
		l := zap.L().With(zap.String("request", logger.MakeRequestID()))
		ctx = logger.SetEntry(ctx, l)

		err := logRequest(l, "Stream "+info.FullMethod, warnDuration, func() error {
			wrapped := grpc_middleware.WrapServerStream(ss)
			wrapped.WrappedContext = ctx
			return handler(srv, ss)
		})
		return err
	}
}
