package servers

import (
	"context"
	"runtime/debug"
	"runtime/pprof"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona-platform/platform/pkg/tracing"

	"github.com/percona-platform/platform/pkg/logger"
)

// logGRPCRequest wraps f (gRPC handler) invocation with logging and panic recovery.
func logGRPCRequest(l *zap.Logger, prefix string, warnD time.Duration, f func() error) (err error) {
	start := time.Now()
	sl := l.Sugar()
	l.Info("Starting " + prefix)

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
		gRPCStatus, isGRPCError := status.FromError(errors.Cause(err))
		switch {
		case err == nil:
			if warnD == 0 || dur < warnD {
				l.Info("Finished "+prefix,
					zap.String("code", gRPCStatus.Code().String()),
					zap.Duration("duration", dur),
				)
			} else {
				l.Warn("Finished "+prefix,
					zap.String("code", gRPCStatus.Code().String()),
					zap.Duration("duration", dur),
					zap.Duration("warn_duration", warnD),
				)
			}
		case isGRPCError:
			l.Warn("Finished "+prefix,
				zap.String("code", gRPCStatus.Code().String()),
				zap.Duration("duration", dur),
				zap.Error(err),
			)
		default:
			l.Error("Finished "+prefix,
				zap.String("code", gRPCStatus.Code().String()),
				zap.Duration("duration", dur),
				zap.Error(err),
			)
			err = status.Error(codes.Internal, "Internal server error.")
		}
	}()

	err = f()

	return //nolint:nakedret
}

// unaryLoggingInterceptor returns a new unary server interceptor that logs incoming requests.
func unaryLoggingInterceptor(l *zap.Logger, warnDuration time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// add pprof labels for more useful profiles
		defer pprof.SetGoroutineLabels(ctx)
		ctx = pprof.WithLabels(ctx, pprof.Labels("method", info.FullMethod))
		pprof.SetGoroutineLabels(ctx)

		// make context with logger
		rl := l
		if reqID := tracing.GetRequestIDFromGrpcIncomingContext(ctx); len(reqID) != 0 {
			rl = rl.With(zap.String("request-id", reqID))
		}
		rl = rl.With(zap.String("method", info.FullMethod))

		// wrap logger into context so that the following gRPC interceptors and handlers could re-use it.
		ctx = logger.GetContextWithLogger(ctx, rl)

		var res interface{}
		err := logGRPCRequest(rl, "unary call", warnDuration, func() error {
			var origErr error
			res, origErr = handler(ctx, req)
			return origErr
		})

		// err is already logged by logRequest
		l.Sugar().Debugf("\nRequest:\n%s\nResponse:\n%s\n", req, res)

		return res, err
	}
}

// streamLoggingInterceptor returns a new stream server interceptor that logs incoming messages.
func streamLoggingInterceptor(l *zap.Logger, warnDuration time.Duration) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		// add pprof labels for more useful profiles
		defer pprof.SetGoroutineLabels(ctx)
		ctx = pprof.WithLabels(ctx, pprof.Labels("method", info.FullMethod))
		pprof.SetGoroutineLabels(ctx)

		// make context with logger
		rl := l
		if requestID := tracing.GetRequestIDFromGrpcIncomingContext(ctx); len(requestID) != 0 {
			rl = rl.With(zap.String("request-id", requestID))
		}
		rl = rl.With(zap.String("method", info.FullMethod))

		// wrap logger into context so that the following gRPC interceptors and handlers could re-use it.
		ctx = logger.GetContextWithLogger(ctx, rl)

		err := logGRPCRequest(l, "stream", warnDuration, func() error {
			wrapped := grpc_middleware.WrapServerStream(ss)
			wrapped.WrappedContext = ctx
			return handler(srv, ss)
		})

		return err
	}
}
