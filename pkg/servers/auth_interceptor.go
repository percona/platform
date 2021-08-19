package servers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	AuthSessionHeader = "Auth-Session" // Okta authentication session ID
	AuthEmailHeader   = "Auth-Email"   // user's email
	AuthStatusHeader  = "Auth-Status"  // gRPC status code (codes.Code)
	AuthErrorHeader   = "Auth-Error"   // gRPC error message, if code is not codes.OK
)

var (
	errInvalidCredentials = status.Error(codes.Unauthenticated, "Invalid credentials.")
	errAuthenticationFail = status.Error(codes.Internal, "Authentication fail.")
)

// CustomeHeaderMatcher lets the Auth-* headers added by forwardauth pass
// through to the grpc server when an HTTP request hits grpc-gateway server.
func CustomHeaderMatcher(key string) (string, bool) {
	switch key {
	case AuthSessionHeader:
	case AuthEmailHeader:
	case AuthStatusHeader:
	case AuthErrorHeader:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
	return runtime.DefaultHeaderMatcher(key)
}

func unaryAuthInterceptor(noAuthMethods []string) grpc.UnaryServerInterceptor {
	noAuthMethodsSet := make(map[string]struct{}, len(noAuthMethods))
	for _, m := range noAuthMethods {
		noAuthMethodsSet[m] = struct{}{}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := logger.Get(ctx).Sugar()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error("No metadata in incoming request.")
			return nil, errAuthenticationFail
		}
		l.Debugf("Received metadata: %+v", md)

		// Check authentication error before checking if methods requires authentication at all:
		// * if Authorization header is absent, Auth Service returns OK;
		// * but if Authorization header is present, it should be valid.
		if err := handleAuthProxyError(md, l); err != nil {
			return nil, err
		}

		email, sessionID, err := getAuthData(md, l)
		if err != nil {
			if _, ok := noAuthMethodsSet[info.FullMethod]; !ok {
				return nil, err
			}
		}

		return handler(rdata.AddToContext(ctx, sessionID, email), req)
	}
}

func streamAuthInterceptor(noAuthMethods []string) grpc.StreamServerInterceptor {
	noAuthMethodsSet := make(map[string]struct{}, len(noAuthMethods))
	for _, m := range noAuthMethods {
		noAuthMethodsSet[m] = struct{}{}
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		l := logger.Get(ctx).Sugar()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error("No metadata in incoming request.")
			return errAuthenticationFail
		}
		l.Debugf("Received metadata: %+v", md)

		// Check authentication error before checking if methods requires authentication at all:
		// * if Authorization header is absent, Auth Service returns OK;
		// * but if Authorization header is present, it should be valid.
		if err := handleAuthProxyError(md, l); err != nil {
			return err
		}

		email, sessionID, err := getAuthData(md, l)
		if err != nil {
			if _, ok := noAuthMethodsSet[info.FullMethod]; !ok {
				return err
			}
		}

		return handler(rdata.AddToContext(ctx, sessionID, email), ss)
	}
}

// handleAuthProxyError checks authentication status and message forwarded from proxy
// and returns proper response to user in case on any problem.
func handleAuthProxyError(md metadata.MD, l *zap.SugaredLogger) error {
	authStatus, err := getAuthStatusFromMetadata(md)
	if err != nil {
		l.Errorf("failed to get auth status from request metadata, reason: %+v", err)
		return errAuthenticationFail
	}

	if authStatus != codes.OK {
		authError, err := getAuthErrorFromMetadata(md)
		if err != nil {
			l.Error(err)
		}
		return status.Error(authStatus, authError)
	}

	return nil
}

// TODO Merge five functions below and some code above into function that parses incoming headers/metadata
// and returns struct with four fields. Use rdata package there?

// getAuthData extracts user email and session id from request metadata.
func getAuthData(md metadata.MD, l *zap.SugaredLogger) (string, string, error) {
	email, err := getAuthEmailFromMetadata(md)
	if err != nil {
		l.Errorf("failed to get auth email from request metadata, reason: %+v", err)
		return "", "", errAuthenticationFail
	}

	if email == "" {
		return "", "", errInvalidCredentials
	}

	sessionID, err := getAuthSessionIDFromMetadata(md)
	if err != nil {
		l.Errorf("failed to get auth session id from request metadata, reason: %+v", err)
		return "", "", errAuthenticationFail
	}

	return email, sessionID, nil
}

// getAuthStatusFromMetadata extracts auth status set by proxy from metadata.
func getAuthStatusFromMetadata(md metadata.MD) (codes.Code, error) {
	header := md.Get(AuthStatusHeader)
	if len(header) != 1 {
		return 0, fmt.Errorf("expect exactly one auth status header, got: %d", len(header))
	}

	c, err := strconv.Atoi(header[0])
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse auth status code")
	}

	return codes.Code(c), nil
}

// getAuthErrorFromMetadata extracts auth error message set by proxy from metadata.
func getAuthErrorFromMetadata(md metadata.MD) (string, error) {
	header := md.Get(AuthErrorHeader)
	if len(header) != 1 {
		return "", fmt.Errorf("expect exactly one auth error header, got: %d", len(header))
	}

	return header[0], nil
}

// getAuthEmailFromMetadata extracts user email set by proxy from metadata.
func getAuthEmailFromMetadata(md metadata.MD) (string, error) {
	header := md.Get(AuthEmailHeader)
	if len(header) > 1 {
		return "", fmt.Errorf("expect at most one auth email header, got: %d", len(header))
	}

	if len(header) == 0 {
		return "", nil
	}

	return header[0], nil
}

// getAuthSessionIDFromMetadata extracts user session id set by proxy from metadata.
func getAuthSessionIDFromMetadata(md metadata.MD) (string, error) {
	header := md.Get(AuthSessionHeader)
	if len(header) > 1 {
		return "", fmt.Errorf("expect at most one auth session header, got: %d", len(header))
	}

	if len(header) == 0 {
		return "", nil
	}

	return header[0], nil
}
