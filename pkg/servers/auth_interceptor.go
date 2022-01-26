package servers

import (
	"context"
	"fmt"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/percona-platform/platform/pkg/logger"
	"github.com/percona-platform/platform/pkg/rdata"
	"github.com/percona-platform/platform/pkg/tracing"
)

type authMethodType int

// Headers set by proxy.
const (
	// AuthUsernameHeader Percona Account username that is used for authentication.
	AuthUsernameHeader = "Auth-Username"

	// AuthUserIDHeader Percona Account User ID in Okta.
	// Note: Percona Account is handled by Okta so ID comes from Okta as well.
	AuthUserIDHeader = "Auth-User-ID"

	// AuthAppIDHeader Application ID in Okta.
	// Note: Application is handled by Okta so ID comes from Okta as well.
	AuthAppIDHeader = "Auth-App-ID"

	// AuthSuperAdminHeader flag indicates that this particular user has SuperAdmin
	// permissions in Percona Portal only.
	AuthSuperAdminHeader = "Auth-Portal-Super-Admin"

	// AuthPortalOrgIDHeader Percona Portal Organization ID (equal to Okta Group ID).
	AuthPortalOrgIDHeader = "Auth-Portal-Org-ID"

	// AuthTokenHeader holds OAuth2 access_token that was used for request authentication.
	// Is used for token propagation to outgoing requests since 'Authorization'
	// HTTP header is removed by Traefik after request authentication.
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

	// Auth methods types.

	// Method doesn't require authentication.
	authMethodNoAuth authMethodType = iota

	// Method may use authentication data if it exists.
	// Anonymous calls of this method are allowed as well.
	authMethodMayUseAuth

	// Method require authentication otherwise it will be rejected.
	authMethodRequireAuth
)

var errAuthenticationFail = status.Error(codes.Unauthenticated, "Authentication fail.")

// PerconaHeaderMatcher preserves the Auth-* headers added by /forwardauth in Authed service
// after the HTTP request is received by grpc-gateway and are forwarded as-is
// to the grpc server.
// It also preserves tracing headers.
func PerconaHeaderMatcher(key string) (string, bool) {
	keyCanonical := textproto.CanonicalMIMEHeaderKey(key)
	if perconaAuthHeadersMatcher(keyCanonical) {
		return key, true
	}

	if tracing.OpenTracingHeadersMatcher(keyCanonical) {
		return key, true
	}

	return runtime.DefaultHeaderMatcher(key)
}

// perconaAuthHeadersMatcher filter function for the Percona Auth-* headers added by /forwardauth in Authed service.
// NOTE: key parameter must be in a Canonical format.
func perconaAuthHeadersMatcher(key string) bool {
	return strings.HasPrefix(key, "Auth-")
}

func unaryAuthInterceptor(noAuthMethods, mayUseAuthMethods []string) grpc.UnaryServerInterceptor { //nolint:cyclop, funlen
	noAuthMethodsSet := make(map[string]struct{}, len(noAuthMethods))
	mayUseAuthMethodsSet := make(map[string]struct{}, len(mayUseAuthMethods))

	for _, m := range noAuthMethods {
		noAuthMethodsSet[m] = struct{}{}
	}

	for _, m := range mayUseAuthMethods {
		if _, ok := noAuthMethodsSet[m]; ok {
			panic(fmt.Sprintf("method %s can't be listed in NoAuthMethods and MayUseAuthMethods simultaneously", m))
		}
		mayUseAuthMethodsSet[m] = struct{}{}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := logger.GetLoggerFromContext(ctx)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error("No metadata in incoming request. Rejecting request.")
			return nil, errAuthenticationFail
		}
		l.Debug("Received metadata", zap.Any("metadata", md))

		// Check authentication error before checking if methods requires authentication at all:
		// * if Authorization header is absent, Auth Service returns OK;
		// * but if Authorization header is present, it should be valid.
		if err := handleAuthProxyError(md, l); err != nil {
			l.Error("Incoming request is unauthenticated. Rejecting request.")
			return nil, err
		}

		authData := new(rdata.RequestData)
		var err error

		switch getAuthMethodType(noAuthMethodsSet, mayUseAuthMethodsSet, info.FullMethod) {
		case authMethodRequireAuth:
			// Request must be authenticated.
			authData, err = getAuthData(md)
			if err != nil {
				l.Error("Can't extract auth data from incoming request. Rejecting request.", zap.Error(err))
				return nil, errAuthenticationFail
			}
			ctx = rdata.AddToContext(ctx, authData)
		case authMethodMayUseAuth, authMethodNoAuth:
			// In case auth data exist add it to context.
			tmpAuthData, err := getAuthData(md)
			if err == nil {
				authData = tmpAuthData
				ctx = rdata.AddToContext(ctx, authData)
			}
		default:
			// Do not try to extract auth data from incoming context.
		}

		// Add logger with userID/appID attributes to context.
		// This logger will be extracted from context and used later by service layers.
		zapUserID := zap.Skip()
		if len(authData.UserID) != 0 {
			zapUserID = zap.String(logger.UserIDAttr, authData.UserID)
		}

		zapAppID := zap.Skip()
		if len(authData.AppID) != 0 {
			zapAppID = zap.String(logger.AppIDAttr, authData.AppID)
		}

		ctx = logger.GetContextWithLogger(ctx, l.With(zapUserID, zapAppID))
		return handler(ctx, req)
	}
}

func streamAuthInterceptor(noAuthMethods, mayUseAuthMethods []string) grpc.StreamServerInterceptor { //nolint:cyclop, funlen
	noAuthMethodsSet := make(map[string]struct{}, len(noAuthMethods))
	mayUseAuthMethodsSet := make(map[string]struct{}, len(mayUseAuthMethods))

	for _, m := range noAuthMethods {
		noAuthMethodsSet[m] = struct{}{}
	}

	for _, m := range mayUseAuthMethods {
		if _, ok := noAuthMethodsSet[m]; ok {
			panic(fmt.Sprintf("method %s can't be listed in NoAuthMethods and MayUseAuthMethods simultaneously", m))
		}
		mayUseAuthMethodsSet[m] = struct{}{}
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		l := logger.GetLoggerFromContext(ctx)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error("No metadata in incoming request. Rejecting request.")
			return errAuthenticationFail
		}
		l.Debug("Received metadata", zap.Any("metadata", md))

		// Check authentication error before checking if methods requires authentication at all:
		// * if Authorization header is absent, Auth Service returns OK;
		// * but if Authorization header is present, it should be valid.
		if err := handleAuthProxyError(md, l); err != nil {
			l.Error("Incoming request is unauthenticated. Rejecting request.")
			return err
		}

		authData := new(rdata.RequestData)
		var err error

		switch getAuthMethodType(noAuthMethodsSet, mayUseAuthMethodsSet, info.FullMethod) {
		case authMethodRequireAuth:
			// Request must be authenticated.
			authData, err = getAuthData(md)
			if err != nil {
				l.Error("Can't extract auth data from incoming request. Rejecting request.", zap.Error(err))
				return errAuthenticationFail
			}
			ctx = rdata.AddToContext(ctx, authData)
		case authMethodMayUseAuth, authMethodNoAuth:
			// In case auth data exist add it to context.
			tmpAuthData, err := getAuthData(md)
			if err == nil {
				authData = tmpAuthData
				ctx = rdata.AddToContext(ctx, authData)
			}
		default:
			// Do not try to extract auth data from incoming context.
		}

		// Add logger with userID/appID attributes to context.
		// This logger will be extracted from context and used later by service layers.
		zapUserID := zap.Skip()
		if len(authData.UserID) != 0 {
			zapUserID = zap.String(logger.UserIDAttr, authData.UserID)
		}

		zapAppID := zap.Skip()
		if len(authData.AppID) != 0 {
			zapAppID = zap.String(logger.AppIDAttr, authData.AppID)
		}

		ctx = logger.GetContextWithLogger(ctx, l.With(zapUserID, zapAppID))
		return handler(ctx, ss)
	}
}

func getAuthMethodType(noAuthMethodsSet, mayUseAuthMethodsSet map[string]struct{}, m string) authMethodType {
	if _, ok := noAuthMethodsSet[m]; ok {
		return authMethodNoAuth
	}

	if _, ok := mayUseAuthMethodsSet[m]; ok {
		return authMethodMayUseAuth
	}
	return authMethodRequireAuth
}

// handleAuthProxyError checks authentication status and message forwarded from proxy
// and returns proper response to user in case of any problem.
func handleAuthProxyError(md metadata.MD, l *zap.Logger) error {
	authStatus, err := getAuthStatusFromMetadata(md)
	if err != nil {
		l.Error("Failed to get auth status from request metadata.", zap.Error(err))
		return errAuthenticationFail
	}

	if authStatus != codes.OK {
		authError, err := getAuthErrorFromMetadata(md)
		if err != nil {
			l.Error("Failed to extract auth error from incoming request.", zap.Error(err))
			return errAuthenticationFail
		}
		return status.Error(authStatus, authError)
	}

	return nil
}

// getAuthData extracts user email and session id from request metadata.
func getAuthData(md metadata.MD) (*rdata.RequestData, error) { //nolint: funlen, cyclop
	username, err := getStringFromMetadata(md, AuthUsernameHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthUsernameHeader)
	}

	userID, err := getStringFromMetadata(md, AuthUserIDHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthUserIDHeader)
	}

	appID, err := getStringFromMetadata(md, AuthAppIDHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthAppIDHeader)
	}

	isPortalSuperAdmin, err := getBoolFromMetadata(md, AuthSuperAdminHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthSuperAdminHeader)
	}

	portalOrgID, err := getStringFromMetadata(md, AuthPortalOrgIDHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthPortalOrgIDHeader)
	}

	authToken, err := getStringFromMetadata(md, AuthTokenHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthTokenHeader)
	}

	// Keep for backward compatibility.
	email, err := getStringFromMetadata(md, AuthEmailHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthEmailHeader)
	}

	sessionID, err := getStringFromMetadata(md, AuthSessionHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s from request metadata", AuthSessionHeader)
	}

	// There are the following cases possible:
	// - Auth-Username header is not empty - it means this incoming request we are processing now
	// is from real user (browser).
	// - Auth-Portal-Org-ID header is not empty - it means this incoming request we are processing now
	// is from PMM Server (machine-to-machine communication).
	// Authorized incoming request must contain one of: username, portalOrgID, sessionID.
	if len(username) == 0 && len(portalOrgID) == 0 && len(sessionID) == 0 {
		return nil, fmt.Errorf("at least one of the auth headers [%s,%s,%s] must be provided", AuthUsernameHeader, AuthPortalOrgIDHeader, AuthSessionHeader)
	}

	return &rdata.RequestData{
		Username:           username,
		UserID:             userID,
		AppID:              appID,
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
		return "", fmt.Errorf("expect exactly one %s header, got: %d", AuthErrorHeader, len(header))
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
