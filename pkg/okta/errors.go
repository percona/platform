package okta

import "github.com/okta/okta-sdk-golang/v2/okta"

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

// Error returns error message. If error cause is Okta error it will add Okta error summary to message.
func (e *AuthError) Error() string {
	if e.origin != nil {
		if oErr, ok := e.origin.(*okta.Error); ok {
			return e.msg + ": " + getOktaErrorCause(oErr)
		}
		return e.msg + ": " + e.origin.Error()
	}
	return e.msg
}

// OriginError returns origin error that causes AuthError if exists.
func (e *AuthError) OriginError() error {
	return e.origin
}

// getOktaErrorCause extracts cause message from Okta error.
func getOktaErrorCause(e *okta.Error) string {
	if len(e.ErrorCauses) == 0 {
		return ""
	}

	cause, ok := e.ErrorCauses[0]["errorSummary"]
	if !ok {
		return ""
	}

	return cause.(string)
}
