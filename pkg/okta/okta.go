// Package okta implements methods for interacting with Okta API.
package okta

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/logger"
)

// Client implements methods for interacting with Okta Identity Service API.
// Some methods can return AuthError which indicates on authentication/authorisation problems.
type Client struct {
	l            *zap.SugaredLogger
	c            *okta.Client
	oktaHost     string
	oktaAPIToken string
}

// New returns new Service instance.
func New(ctx context.Context, host, token string) (*Client, error) {
	l := zap.L().Named("okta").Sugar()

	u := url.URL{Scheme: "https", Host: host}

	_, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(u.String()),
		okta.WithToken(token),
		okta.WithHttpClientPtr(&http.Client{
			Transport: logger.HTTP(http.DefaultTransport, l.Debugf),
		}),
		okta.WithCache(false),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		l:            l,
		c:            client,
		oktaHost:     host,
		oktaAPIToken: token,
	}, nil
}

// SignUp creates user account with provided login, first name and last name.
// The user is created with a "PROVISIONED" state and a verification email is sent to the user.
// Once the user sets the password the state changes to "ACTIVE" in Okta.
// Returns AuthError when login, firstName and lastName violates validation rules, also when login already exists.
func (c *Client) SignUp(ctx context.Context, login, firstName, lastName string) (*User, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}

	if firstName == "" {
		return nil, ErrEmptyFirstName
	}

	if lastName == "" {
		return nil, ErrEmptyLastName
	}

	u := okta.CreateUserRequest{
		Profile: &okta.UserProfile{
			"email":     login,
			"login":     login,
			"firstName": firstName,
			"lastName":  lastName,
		},
	}
	qp := query.NewQueryParams(query.WithActivate(true))
	user, _, err := c.c.User.CreateUser(ctx, u, qp)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrap(err, "failed to sign up user")
	}

	nLogin, err := getUserLogin(user)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", user.Id)
	}

	return &User{
		ID:     user.Id,
		Login:  nLogin,
		Status: user.Status,
	}, nil
}

// FindUser searches user by login and returns user.
func (c *Client) FindUser(ctx context.Context, login string) (*User, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	user, _, err := c.c.User.GetUser(ctx, url.QueryEscape(login))
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrapf(err, "failed to find user")
	}

	userLogin, err := getUserLogin(user)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", user.Id)
	}

	firstName, err := getUserFirstName(user)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", user.Id)
	}

	lastName, err := getUserLastName(user)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", user.Id)
	}

	return &User{
		ID:        user.Id,
		Login:     userLogin,
		FirstName: firstName,
		LastName:  lastName,
		Status:    user.Status,
	}, nil
}

// SignIn returns user id and session oktaAPIToken. Returns AuthError in case of invalid login or password.
func (c *Client) SignIn(ctx context.Context, login, password string) (string, string, error) {
	if password == "" {
		return "", "", ErrEmptyPassword
	}

	if login == "" {
		return "", "", ErrEmptyLogin
	}

	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: login,
		Password: password,
	}

	resp := struct {
		ExpiresAt    string `json:"expiresAt"`
		Status       string `json:"status"`
		SessionToken string `json:"sessionToken"`
		Embedded     struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
		} `json:"_embedded"` //nolint:tagliatelle
	}{}

	err := c.DoRequest(ctx, "POST", "/api/v1/authn", data, &resp)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return "", "", convertOktaError(oErr)
		}

		return "", "", errors.Wrap(err, "failed to authenticate user")
	}

	return resp.Embedded.User.ID, resp.SessionToken, nil
}

// DeleteUser deactivates and deletes user.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	err := c.deactivateUser(ctx, userID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if err == nil {
		nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := c.waitForDeactivation(nCtx, userID); err != nil {
			return err
		}
	}

	_, err = c.c.User.DeactivateOrDeleteUser(ctx, userID, nil)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}

// UpdateProfile updates user's profile.
func (c *Client) UpdateProfile(ctx context.Context, user *User, firstName, lastName string) (*User, error) {
	// The okta-go-sdk implementation differs slightly from the API docs where it uses
	// PUT instead of POST for this endpoint. In case of PUT (used by the SDK) it requires
	// both the email and login fields to be non-empty, so we send the original login to prevent
	// the request from returning an error.
	body := okta.User{ //nolint:exhaustivestruct
		Profile: &okta.UserProfile{
			"firstName": firstName,
			"lastName":  lastName,
			"login":     user.Login,
			"email":     user.Login,
		},
	}
	u, _, err := c.c.User.UpdateUser(ctx, user.ID, body, nil)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrapf(err, "failed to update user profile")
	}

	userLogin, err := getUserLogin(u)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", u.Id)
	}

	fName, err := getUserFirstName(u)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", u.Id)
	}

	lName, err := getUserLastName(u)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", u.Id)
	}

	return &User{
		ID:        u.Id,
		Login:     userLogin,
		FirstName: fName,
		LastName:  lName,
		Status:    u.Status,
	}, nil
}

func (c *Client) deactivateUser(ctx context.Context, userID string) error {
	_, err := c.c.User.DeactivateUser(ctx, userID, nil)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to deactivate user")
	}

	return nil
}

func (c *Client) waitForDeactivation(ctx context.Context, userID string) error {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		user, _, err := c.c.User.GetUser(ctx, userID)
		if err != nil {
			var oErr *okta.Error
			if errors.As(err, &oErr) {
				return convertOktaError(oErr)
			}

			return errors.Wrapf(err, "failed to check user status")
		}

		if user.Status == "DEPROVISIONED" {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			// next loop iteration
		}
	}
}

// GetRegisteredUsersCount returns number of regustered users.
func (c *Client) GetRegisteredUsersCount(ctx context.Context) (float64, error) {
	qp := query.NewQueryParams(query.WithQ("Everyone"), query.WithFilter("type eq \"BUILT_IN\""))
	groups, _, err := c.c.Group.ListGroups(ctx, qp)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return 0, convertOktaError(oErr)
		}

		return 0, errors.Wrap(err, "failed to find everyone users group")
	}

	if len(groups) != 1 {
		return 0, fmt.Errorf("expect only one everyone group search result, got %d", len(groups))
	}

	var group okta.Group
	path := fmt.Sprintf("/api/v1/groups/%s?expand=stats", groups[0].Id)
	err = c.DoRequest(ctx, "GET", path, nil, &group)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return 0, convertOktaError(oErr)
		}

		return 0, errors.Wrap(err, "failed to get group stats")
	}

	embedded, ok := (group.Embedded).(map[string]interface{})
	if !ok {
		return 0, errors.New("missing embedded section in group stats response")
	}

	stats, ok := embedded["stats"].(map[string]interface{})
	if !ok {
		return 0, errors.New("missing stats section in group stats response")
	}

	value, ok := stats["usersCount"]
	if !ok {
		return 0, errors.New("missing usersCount section in group stats response")
	}

	usersCount, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("can't cast usersCount to float64, unexpected type %T", usersCount)
	}

	return usersCount, nil
}

// CreateSession creates session and returns session id and expiration date.
func (c *Client) CreateSession(ctx context.Context, sessionToken string) (string, time.Time, error) {
	session, _, err := c.c.Session.CreateSession(ctx, okta.CreateSessionRequest{SessionToken: sessionToken})
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to create session")
	}

	return session.Id, *session.ExpiresAt, nil
}

// CheckSession returns user login if session exists and active.
func (c *Client) CheckSession(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", ErrNotFound
	}

	session, _, err := c.c.Session.GetSession(ctx, sessionID)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return "", convertOktaError(oErr)
		}

		return "", errors.Wrap(err, "failed to check session")
	}

	return session.Login, nil
}

// RefreshSession resets session timeout and returns new expiration date.
func (c *Client) RefreshSession(ctx context.Context, sessionID string) (time.Time, error) {
	session, _, err := c.c.Session.RefreshSession(ctx, sessionID)
	if err != nil || session.ExpiresAt == nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return time.Time{}, convertOktaError(oErr)
		}

		return time.Time{}, errors.Wrap(err, "failed to refresh session")
	}

	return *session.ExpiresAt, nil
}

// CloseSession terminates user session.
func (c *Client) CloseSession(ctx context.Context, sessionID string) error {
	_, err := c.c.Session.EndSession(ctx, sessionID)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to close session")
	}

	return nil
}

// CreateGroup creates group with provided name and description.
func (c *Client) CreateGroup(ctx context.Context, name, description string) (*Group, error) {
	req := okta.Group{
		Profile: &okta.GroupProfile{
			Name:        name,
			Description: description,
		},
	}

	group, _, err := c.c.Group.CreateGroup(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create group")
	}

	return &Group{
		ID:          group.Id,
		Name:        group.Profile.Name,
		Description: group.Profile.Description,
	}, nil
}

// GroupExists finds whether okta group with the provided name exists.
func (c *Client) GroupExists(ctx context.Context, name string) (bool, error) {
	var g []okta.Group
	params := url.Values{}
	params.Add("q", name)
	params.Add("limit", "1")

	err := c.DoRequest(ctx, "GET", fmt.Sprintf("/api/v1/groups?%s", params.Encode()), nil, &g)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return false, convertOktaError(oErr)
		}

		return false, errors.Wrap(err, "failed to find group")
	}

	for _, group := range g {
		// Okta matches groups in case-insensitive format.
		if strings.EqualFold(group.Profile.Name, name) {
			return true, nil
		}
	}
	return false, nil
}

// DeleteGroup delete group with provided ID.
func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := c.c.Group.DeleteGroup(ctx, groupID)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to delete group")
	}

	return nil
}

// GetGroupMembers returns list of group members.
func (c *Client) GetGroupMembers(ctx context.Context, groupID string, limit int, cursor string) ([]User, error) {
	params := query.Params{
		Limit:  int64(limit),
		Cursor: cursor,
	}
	users, _, err := c.c.Group.ListGroupUsers(ctx, groupID, &params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get group members")
	}

	res := make([]User, 0, len(users))
	for _, user := range users {
		login, err := getUserLogin(user)
		if err != nil {
			c.l.Warnf("User %s has bad profile, reason: %+v.", user.Id, err)
		}

		res = append(res, User{ID: user.Id, Login: login, Status: user.Status})
	}

	return res, nil
}

// AddUserToGroup add user to group.
func (c *Client) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	_, err := c.c.Group.AddUserToGroup(ctx, groupID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user to group")
	}

	return nil
}

// IsAppAssignedToGroup returns true if app is assigned to given group, false otherwise - not found error. Also false in case of any other error.
func (c *Client) IsAppAssignedToGroup(ctx context.Context, appID, groupID string) bool {
	if err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/apps/%s/groups/%s", appID, groupID), nil, nil); err != nil {
		return false
	}
	return true
}

// AddAppToGroup adds app to group.
func (c *Client) AddAppToGroup(ctx context.Context, appID, groupID string) error {
	err := c.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/apps/%s/groups/%s", appID, groupID), nil, nil)
	if err != nil {
		return errors.Wrap(err, "failed to add app to the group")
	}

	return nil
}

// RemoveAppFromGroup removes app from group.
func (c *Client) RemoveAppFromGroup(ctx context.Context, appID, groupID string) error {
	err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/apps/%s/groups/%s", appID, groupID), nil, nil)
	if err != nil {
		return errors.Wrap(err, "failed to remove app from the group")
	}

	return nil
}

// ErrOriginNotFound means operation on origin failed because it does not exist.
var ErrOriginNotFound error = errors.New("trusted origin was not found")

// GetTrustedOriginID returns origin's id if it exists, nil and error when it does not.
func (c *Client) GetTrustedOriginID(ctx context.Context, origin string) (string, error) {
	origins, response, err := c.c.TrustedOrigin.ListOrigins(ctx, nil)
	if response.HasNextPage() {
		c.l.Warn("The list of origins is not complete. The trusted origins API got support for pagination!")
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to check if origin is trusted")
	}
	for _, trusted := range origins {
		if trusted.Origin == origin {
			return trusted.Id, nil
		}
	}
	return "", ErrOriginNotFound
}

// AddTrustedOrigin makes the given origin trusted and returns it's id.
func (c *Client) AddTrustedOrigin(ctx context.Context, origin string) (string, error) {
	trusted, _, err := c.c.TrustedOrigin.CreateOrigin(ctx, okta.TrustedOrigin{
		Name:   origin,
		Origin: origin,
		Scopes: []*okta.Scope{
			{Type: "REDIRECT"},
		},
	})
	if err != nil {
		return "", err
	}
	return trusted.Id, nil
}

// DeleteTrustedOrigin deletes the given trusted origin.
func (c *Client) DeleteTrustedOrigin(ctx context.Context, originID string) error {
	_, err := c.c.TrustedOrigin.DeleteOrigin(ctx, originID)
	return err
}

// RemoveUserFromGroup remove user from group.
func (c *Client) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	_, err := c.c.Group.RemoveUserFromGroup(ctx, groupID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user to group")
	}

	return nil
}

// ResetPassword resets user password and sends email for setting new one.
func (c *Client) ResetPassword(ctx context.Context, userID string) error {
	qp := query.NewQueryParams(query.WithSendEmail(true))
	_, _, err := c.c.User.ResetPassword(ctx, userID, qp)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to reset password")
	}

	return nil
}

// UpdateSchema updates schema for provided type.
func (c *Client) UpdateSchema(ctx context.Context, typeID string, schema *Schema) (*Schema, error) {
	var res Schema
	err := c.DoRequest(ctx, "POST", "/api/v1/meta/schemas/user/"+typeID, schema, &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update schema")
	}

	return &res, nil
}

// GetSchema returns schema for provided type.
func (c *Client) GetSchema(ctx context.Context, typeID string) (*Schema, error) {
	var schema Schema
	err := c.DoRequest(ctx, "GET", "/api/v1/meta/schemas/user/"+typeID, nil, &schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schema")
	}

	return &schema, nil
}

func (c *Client) ListPolicies(ctx context.Context, qp *query.Params) ([]*okta.Policy, error) {
	policies, _, err := c.c.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return nil, err
	}
	return policies, err
}

const createOAuthAppRequestBody = `
{
    "name": "oidc_client",
    "label": "PMM-{{ .PMMServerID }}",
    "status": "ACTIVE",
    "signOnMode": "OPENID_CONNECT",
    "credentials": {
        "oauthClient": {
            "autoKeyRotation": true,
            "token_endpoint_auth_method": "client_secret_basic"
        }
    },
    "settings": {
        "oauthClient": {
            "redirect_uris": [
                "{{ .PMMServerCallbackURL }}"
            ],
            "post_logout_redirect_uris": [
                "{{ .PMMServerURL }}"
            ],
            "response_types": [
                "code"
            ],
            "grant_types": [
                "client_credentials",
                "authorization_code",
                "refresh_token"
            ],
            "application_type": "web",
            "consent_method": "REQUIRED",
            "issuer_mode": "CUSTOM_URL"
        }
    },
    "profile": {
        "percona": {
            "portal": {
                "orgId": "{{ .OrgID }}",
                "invId": "{{ .InventoryID }}"
            }
        }
    }
}`

var createOAuthAppRequestBodyTmpl = template.Must(template.New("CreateOAuthAppRequest").Parse(createOAuthAppRequestBody))

// OAuthApp represents an oauth app.
type OAuthApp struct {
	AppID       string `json:"id"`
	Credentials struct {
		OAuthClient struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"oauthClient"`
	} `json:"credentials"`
}

// OAuthAppParams contains values needed when creating a new OAuth app.
type OAuthAppParams struct {
	PMMServerID          string
	PMMServerURL         string
	PMMServerCallbackURL string
	OrgID                string
	InventoryID          string
}

// CreateOAuthApp creates a new OAuth app.
func (c *Client) CreateOAuthApp(ctx context.Context, params *OAuthAppParams) (*OAuthApp, error) {
	var request bytes.Buffer
	err := createOAuthAppRequestBodyTmpl.Execute(&request, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request body for adding the OAuth App")
	}

	var result OAuthApp
	err = c.DoRequest(ctx, "POST", "/api/v1/apps", &request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "adding OAuth APP to Okta failed")
	}
	return &result, nil
}

// DeleteApp deletes an app with given appID.
func (c *Client) DeleteApp(ctx context.Context, appID string) error {
	_, err := c.c.Application.DeactivateApplication(ctx, appID)
	if err != nil {
		return errors.Wrap(err, "failed to deactivate app before deleting it")
	}

	if _, err = c.c.Application.DeleteApplication(ctx, appID); err != nil {
		_, e := c.c.Application.ActivateApplication(ctx, appID)
		if e != nil {
			c.l.Error("Failed to re-activate app after deleting failed, manual intervention in Okta required")
		}
		return err
	}
	return nil
}

// DoRequest makes HTTP requests to okta endpoints.
func (c *Client) DoRequest(ctx context.Context, method, path string, body, v interface{}) error {
	requestExecutor := c.c.CloneRequestExecutor().WithAccept("application/json").WithContentType("application/json")
	req, err := requestExecutor.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	resp, err := requestExecutor.Do(ctx, req, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	return err
}

func getUserLogin(user *okta.User) (string, error) {
	if user.Profile == nil {
		return "", errors.New("missing user profile")
	}

	profile := *user.Profile
	login, ok := profile["login"]
	if !ok {
		return "", errors.New("missing user login")
	}

	return login.(string), nil
}

func getUserFirstName(user *okta.User) (string, error) {
	if user.Profile == nil {
		return "", errors.New("missing user profile")
	}

	profile := *user.Profile
	name, ok := profile["firstName"]
	if !ok {
		return "", errors.New("missing user firstName")
	}

	return name.(string), nil
}

func getUserLastName(user *okta.User) (string, error) {
	if user.Profile == nil {
		return "", errors.New("missing user profile")
	}

	profile := *user.Profile
	name, ok := profile["lastName"]
	if !ok {
		return "", errors.New("missing user lastName")
	}

	return name.(string), nil
}

func convertOktaError(err *okta.Error) error {
	switch err.ErrorCode {
	case "E0000001":
		switch err.ErrorSummary {
		case "Api validation failed: password":
			return NewError("invalid password", err)
		case "Api validation failed: login":
			return NewError("invalid login", err)
		default:
			return err
		}
	case "E0000004":
		return ErrAuthentication
	case "E0000007":
		return ErrNotFound
	default:
		return err
	}
}
