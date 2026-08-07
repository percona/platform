package servers

import (
	"net/textproto"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/percona/platform/pkg/rdata"
	"github.com/percona/platform/pkg/tracing"
)

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

	// AuthHook indicates if the request authorized as hook request.
	AuthHook = "Auth-Hook"

	// OktaVerificationHeader header used by Okta to verify hook handlers.
	OktaVerificationHeader = "X-Okta-Verification-Challenge"

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

// AuthMetadata returns auth headers.
func AuthMetadata(r *rdata.RequestData) metadata.MD {
	return metadata.Pairs(
		AuthStatusHeader, strconv.Itoa(int(codes.OK)),
		AuthTokenHeader, r.AuthToken,
		AuthUsernameHeader, r.Username,
		AuthSessionHeader, r.SessionID,
		AuthHook, strconv.FormatBool(r.Hook),
	)
}

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

	if keyCanonical == OktaVerificationHeader {
		return key, true
	}

	return runtime.DefaultHeaderMatcher(key)
}

// perconaAuthHeadersMatcher filter function for the Percona Auth-* headers added by /forwardauth in Authed service.
// NOTE: key parameter must be in a Canonical format.
func perconaAuthHeadersMatcher(key string) bool {
	return strings.HasPrefix(key, "Auth-")
}
