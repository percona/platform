package servers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/percona-platform/platform/pkg/logger"
	"github.com/percona-platform/platform/pkg/rdata"
)

// Headers set by proxy.
const (
	AuthSessionHeader = "Auth-Session"
	AuthEmailHeader   = "Auth-Email"
	AuthStatusHeader  = "Auth-Status"
	AuthErrorHeader   = "Auth-Error"
)

var (
	errInvalidCredentials = status.Error(codes.Unauthenticated, "Invalid credentials.")
	errAuthenticationFail = status.Error(codes.Internal, "Authentication fail.")
)

func unaryAuthInterceptor(noAuthMethods []string) grpc.UnaryServerInterceptor {
	m := make(map[string]struct{}, len(noAuthMethods))
	for _, method := range noAuthMethods {
		m[method] = struct{}{}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := logger.Get(ctx).Sugar()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errInvalidCredentials
		}

		if err := getAuthStatus(md, l); err != nil {
			return nil, err
		}

		if _, ok := m[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		email, sessionID, err := getAuthData(md, l)
		if err != nil {
			return nil, err
		}

		return handler(rdata.AddToContext(ctx, sessionID, email), req)
	}
}

func streamAuthInterceptor(noAuthMethods []string) grpc.StreamServerInterceptor {
	m := make(map[string]struct{}, len(noAuthMethods))
	for _, method := range noAuthMethods {
		m[method] = struct{}{}
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		l := logger.Get(ctx).Sugar()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return errInvalidCredentials
		}

		if err := getAuthStatus(md, l); err != nil {
			return err
		}

		if _, ok := m[info.FullMethod]; ok {
			return handler(ctx, ss)
		}

		email, sessionID, err := getAuthData(md, l)
		if err != nil {
			return err
		}

		return handler(rdata.AddToContext(ctx, sessionID, email), ss)
	}
}

func getAuthStatus(md metadata.MD, l *zap.SugaredLogger) error {
	authStatus, err := getAuthStatusFormMetadata(md)
	if err != nil {
		l.Errorf("failed to get auth status form request metadata, reason: %+v", err)
		return errAuthenticationFail
	}

	if authStatus != codes.OK {
		authError, err := getAuthErrorFormMetadata(md)
		if err != nil {
			l.Error(err)
		}
		return status.Error(authStatus, authError)
	}

	return nil
}

func getAuthData(md metadata.MD, l *zap.SugaredLogger) (string, string, error) {
	email, err := getAuthEmailFromMetadata(md)
	if err != nil {
		l.Errorf("failed to get auth email form request metadata, reason: %+v", err)
		return "", "", errAuthenticationFail
	}

	if email == "" {
		return "", "", errInvalidCredentials
	}

	sessionID, err := getAuthSessionIDFromMetadata(md)
	if err != nil {
		l.Errorf("failed to get auth session id form request metadata, reason: %+v", err)
		return "", "", errAuthenticationFail
	}

	return email, sessionID, nil
}

func getAuthStatusFormMetadata(md metadata.MD) (codes.Code, error) {
	header := md.Get(AuthStatusHeader)
	if len(header) != 1 {
		return 0, fmt.Errorf("expect one auth status header, got: %d", len(header))
	}

	var code codes.Code
	if err := code.UnmarshalJSON([]byte(header[0])); err != nil {
		return 0, errors.Wrap(err, "failed to parse auth status code")
	}

	return code, nil
}

func getAuthErrorFormMetadata(md metadata.MD) (string, error) {
	header := md.Get(AuthErrorHeader)
	if len(header) != 1 {
		return "", fmt.Errorf("expect one or nauth error header, got: %d", len(header))
	}

	return header[0], nil
}

func getAuthEmailFromMetadata(md metadata.MD) (string, error) {
	header := md.Get(AuthEmailHeader)
	if len(header) > 1 {
		return "", fmt.Errorf("expect one or nauth error header, got: %d", len(header))
	}

	if len(header) == 0 {
		return "", nil
	}

	return header[0], nil
}

func getAuthSessionIDFromMetadata(md metadata.MD) (string, error) {
	header := md.Get(AuthSessionHeader)
	if len(header) > 1 {
		return "", fmt.Errorf("expect one or nauth error header, got: %d", len(header))
	}

	if len(header) == 0 {
		return "", nil
	}

	return header[0], nil
}
