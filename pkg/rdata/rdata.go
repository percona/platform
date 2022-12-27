// Package rdata provides functionality for operating request metadata.
package rdata

import (
	"context"

	"github.com/pkg/errors"
)

// ErrNoAuthData helper error for case when auth data is absent in incoming context.
var ErrNoAuthData = errors.New("missing request data in request context")

type requestDataKey struct{}

// RequestData contains request related data added by interceptors.
// See description: https://confluence.percona.com/display/PMM/Single+Sign-On+-+Portal+Integration+with+Okta#SingleSignOnPortalIntegrationwithOkta-AuthorizationInfoHeaders
type RequestData struct {
	// Username Percona Account username that is used for authentication.
	Username string

	// UserID Percona Account User ID in Okta.
	// Note: Percona Account is handled by Okta so ID comes from Okta as well.
	UserID string

	// AppID Application ID in Okta.
	// Note: Application is handled by Okta so ID comes from Okta as well.
	AppID string

	// IsPortalSuperAdmin flag indicates that this particular user has SuperAdmin
	// permissions in Percona Portal only.
	IsPortalSuperAdmin bool

	// PortalOrgID Percona Portal Organization ID (equal to Okta Group ID).
	PortalOrgID string

	// AuthToken holds OAuth2 access_token that was used for request authentication.
	// Is used for token propagation to outgoing requests since 'Authorization'
	// HTTP header is removed by Traefik after request authentication.
	AuthToken string

	// Hook is a flag that shows if the request has Hook authorization
	Hook bool

	// HookVerification is a string sent by okta to verify hook handler
	HookVerification string

	// Keep for backward compatibility
	UserEmail string
	SessionID string
}

// AddToContext adds session id and user email to request context.
func AddToContext(ctx context.Context, data *RequestData) context.Context {
	return context.WithValue(ctx, requestDataKey{}, data)
}

// GetFromContext extracts request data from request context.
func GetFromContext(ctx context.Context) (*RequestData, error) {
	v := ctx.Value(requestDataKey{})
	if v != nil {
		if d, ok := v.(*RequestData); ok {
			return d, nil
		}
	}

	return nil, ErrNoAuthData
}
