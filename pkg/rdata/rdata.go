// Package rdata provides functionality for operating request metadata.
package rdata

import (
	"context"

	"github.com/pkg/errors"
)

type requestDataKey struct{}

// RequestData contains request related data added by interceptors.
type RequestData struct {
	UserEmail string
	SessionID string
}

// AddToContext adds session id and user email to request context.
func AddToContext(ctx context.Context, sessionID, userEmail string) context.Context {
	return context.WithValue(ctx, requestDataKey{}, &RequestData{SessionID: sessionID, UserEmail: userEmail})
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
