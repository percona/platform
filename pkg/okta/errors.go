package okta

import (
	"github.com/pkg/errors"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	ErrEmptyLogin     = &AuthError{msg: "login is empty"}     //nolint:golint
	ErrEmptyFirstName = &AuthError{msg: "firstName is empty"} //nolint:golint
	ErrEmptyLastName  = &AuthError{msg: "lastName is empty"}  //nolint:golint
	ErrEmptyPassword  = &AuthError{msg: "password is empty"}  //nolint:golint
	ErrAuthentication = &AuthError{msg: "authentication error"}
	ErrNotFound       = &AuthError{msg: "not found"}
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
