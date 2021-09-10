// Package rdata provides functionality for operating request metadata.
package rdata

import (
	"context"

	"github.com/pkg/errors"
)

type requestDataKey struct{}

// RequestData contains request related data added by interceptors.
// See description: https://confluence.percona.com/display/PMM/Single+Sign-On+-+Portal+Integration+with+Okta#SingleSignOnPortalIntegrationwithOkta-AuthorizationInfoHeaders
type RequestData struct {
	Username           string
	UserID             string
	IsPortalSuperAdmin bool
	PortalOrgID        string
	AuthToken          string

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

	return nil, errors.New("missing request data in request context")
}
