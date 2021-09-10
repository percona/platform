package servers

import (
	"context"
	"fmt"
	"net/textproto"
	"strconv"
	"strings"

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
	// AuthUsernameHeader is Percona Account username.
	AuthUsernameHeader = "Auth-Username"

	// AuthUserIDHeader is Percona Account ID.
	AuthUserIDHeader = "Auth-User-ID"

	// AuthSuperAdminHeader indicates that this
	// Percona Account has Super Admin permissions on Portal.
	AuthSuperAdminHeader = "Auth-Portal-Super-Admin"

	// AuthPortalOrgIDHeader is Portal Organization ID.
	AuthPortalOrgIDHeader = "Auth-Portal-Org-ID"

	// AuthTokenHeader is OAuth2 access_token.
	AuthTokenHeader = "Auth-Token"

	// Keep for backward compatibility.

	// AuthSessionHeader Okta authentication session ID.
	AuthSessionHeader = "Auth-Session"

	// AuthEmailHeader user's email.
	AuthEmailHeader = "Auth-Email"

	// AuthStatusHeader gRPC status code (codes.Code).
	AuthStatusHeader = "Auth-Status"

	// AuthErrorHeader gRPC error message, if code is not codes.OK.
	AuthErrorHeader = "Auth-Error"
)

var (
	errInvalidCredentials = status.Error(codes.Unauthenticated, "Invalid credentials.")
	errAuthenticationFail = status.Error(codes.Internal, "Authentication fail.")
)

// PerconaAuthHeaderMatcher preserves the PP-Auth-* headers added by /forwardauth in Authed service
// after the HTTP request is received by grpc-gateway and are forwarded as-is
// to the grpc server.
func PerconaAuthHeaderMatcher(key string) (string, bool) {
	keyCanonical := textproto.CanonicalMIMEHeaderKey(key)
	if strings.HasPrefix(keyCanonical, "Auth-") {
		return key, true
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

		if _, ok := noAuthMethodsSet[info.FullMethod]; !ok {
			// Request must be authenticated.
			reqData, err := getAuthData(md, l)
			if err != nil {
				return nil, err
			}
			ctx = rdata.AddToContext(ctx, reqData)
		}

		return handler(ctx, req)
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

		if _, ok := noAuthMethodsSet[info.FullMethod]; !ok {
			// Request must be authenticated.
			reqData, err := getAuthData(md, l)
			if err != nil {
				return err
			}
			ctx = rdata.AddToContext(ctx, reqData)
		}

		return handler(ctx, ss)
	}
}

// handleAuthProxyError checks authentication status and message forwarded from proxy
// and returns proper response to user in case of any problem.
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
func getAuthData(md metadata.MD, l *zap.SugaredLogger) (*rdata.RequestData, error) {
	username, err := getStringFromMetadata(md, AuthUsernameHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthUsernameHeader, err)
		return nil, errAuthenticationFail
	}

	userID, err := getStringFromMetadata(md, AuthUserIDHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthUserIDHeader, err)
		return nil, errAuthenticationFail
	}

	isPortalSuperAdmin, err := getBoolFromMetadata(md, AuthSuperAdminHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthSuperAdminHeader, err)
		return nil, errAuthenticationFail
	}

	portalOrgID, err := getStringFromMetadata(md, AuthPortalOrgIDHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthPortalOrgIDHeader, err)
		return nil, errAuthenticationFail
	}

	authToken, err := getStringFromMetadata(md, AuthTokenHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthTokenHeader, err)
		return nil, errAuthenticationFail
	}

	// Keep for backward compatibility.
	email, err := getStringFromMetadata(md, AuthEmailHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthEmailHeader, err)
		return nil, errAuthenticationFail
	}

	sessionID, err := getStringFromMetadata(md, AuthSessionHeader)
	if err != nil {
		l.Errorf("failed to get %s from request metadata, reason: %+v", AuthSessionHeader, err)
		return nil, errAuthenticationFail
	}

	// There are the following cases possible:
	// - username exists in PP-Auth- headers - it means this incoming request we are processing now
	// is from real user (browser).
	// - portalOrgID exists in PP-Auth- headers - it means this incoming request we are processing now
	// is from PMM Server (machine-to-machine communication).
	// Authorized incoming request must contain one of: username, portalOrgID, sessionID.
	if len(username) == 0 && len(portalOrgID) == 0 && len(sessionID) == 0 {
		l.Errorf("at least one of the auth headers [%s,%s,%s] must be provided", AuthUsernameHeader, AuthPortalOrgIDHeader, AuthSessionHeader)
		return nil, errInvalidCredentials
	}

	return &rdata.RequestData{
		Username:           username,
		UserID:             userID,
		IsPortalSuperAdmin: isPortalSuperAdmin,
		PortalOrgID:        portalOrgID,
		AuthToken:          authToken,
		UserEmail:          email,
		SessionID:          sessionID,
	}, nil
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

// getStringFromMetadata extracts string key set by proxy from metadata.
func getStringFromMetadata(md metadata.MD, key string) (string, error) {
	header := md.Get(key)
	if len(header) > 1 {
		return "", fmt.Errorf("expect at most one %s header, got: %d", key, len(header))
	}

	if len(header) == 0 {
		return "", nil
	}

	return header[0], nil
}

// getBoolFromMetadata extracts bool key set by proxy from metadata.
func getBoolFromMetadata(md metadata.MD, key string) (bool, error) {
	header := md.Get(key)
	if len(header) > 1 {
		return false, fmt.Errorf("expect at most one %s header, got: %d", key, len(header))
	}

	if len(header) == 0 {
		return false, nil
	}

	v, err := strconv.ParseBool(header[0])
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse %s header", key)
	}

	return v, nil
}
