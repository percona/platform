package okta

import (
	"github.com/pkg/errors"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	ErrEmptyLogin                = &AuthError{msg: "login is empty"} //nolint:revive
	ErrEmptyFirstName            = &AuthError{msg: "firstName is empty"}
	ErrEmptyLastName             = &AuthError{msg: "lastName is empty"}
	ErrEmptyPassword             = &AuthError{msg: "password is empty"}
	ErrAuthentication            = &AuthError{msg: "authentication error"}
	ErrNotFound                  = &AuthError{msg: "not found"}
	ErrInvalidPortalAdminOrgs    = &AuthError{msg: "portalAdminOrgs contains invalid uuid"}
	ErrDuplicatedPortalAdminOrgs = &AuthError{msg: "portalAdminOrgs contains duplicated values"}
	ErrEmptyPortalAdminOrgs      = &AuthError{msg: "portalAdminOrgs field is nil"}
)

// AuthError represents authentication/authorisation errors. It contains message that describes
// reason and could contain origin error.
type AuthError struct {
	origin error
	msg    string
}

// NewError returns new AuthError with content.
func NewError(msg string, origin error) error {
	return &AuthError{
		msg:    msg,
		origin: origin,
	}
}

// Error returns error message. If error cause is Okta error it will add Okta error summary to message.
func (e *AuthError) Error() string {
	if e.origin != nil {
		var oErr *okta.Error
		if errors.As(e.origin, &oErr) {
			return e.msg + ": " + oErr.ErrorSummary
		}
		return e.msg + ": " + e.origin.Error()
	}
	return e.msg
}

// Unwrap returns origin error that causes AuthError if exists.
func (e *AuthError) Unwrap() error {
	return e.origin
}
