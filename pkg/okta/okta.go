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

	"github.com/google/uuid"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/logger"
)

// Client implements methods for interacting with Okta Identity Service API.
// Some methods can return AuthError which indicates on authentication/authorisation problems.
type Client struct {
	c                   *okta.Client
	oktaHost            string
	oktaAPIToken        string
	oktaEveryoneGroupID string
}

const (
	profileLastName        = "lastName"
	profileFirstName       = "firstName"
	profileEmail           = "email"
	profileLogin           = "login"
	profilePortalAdminOrgs = "portalAdminOrgs"
)

// New returns new Service instance.
func New(ctx context.Context, host, token string) (*Client, error) {
	u := url.URL{Scheme: "https", Host: host}

	_, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(u.String()),
		okta.WithToken(token),
		okta.WithHttpClientPtr(&http.Client{
			Transport: logger.HTTP(http.DefaultTransport, "oktaClient"),
		}),
		okta.WithCache(false),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
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
	l := extractLogger(ctx)
	l.Info("Creating new Okta user.", zap.String("login", login))

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
			profileLogin:     login,
			profileEmail:     login,
			profileFirstName: firstName,
			profileLastName:  lastName,
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

// FindUser searches user either by login or okta user ID and returns user.
func (c *Client) FindUser(ctx context.Context, login string) (*User, error) {
	l := extractLogger(ctx)
	l.Info("Looking for Okta user by username.", zap.String("username", login))

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

	return convertUser(user)
}

// RegisterUser invites okta user and returns user.
func (c *Client) RegisterUser(ctx context.Context, params RegisterUserParams) (*User, error) {
	l := extractLogger(ctx)

	err := validateRegisterUserParams(params)
	if err != nil {
		return nil, err
	}

	l.Info("Inviting Okta user.", zap.String("login", params.Login))

	activate := false
	profile := okta.UserProfile{
		profileLogin:           params.Login,
		profileEmail:           params.Login,
		profileFirstName:       "",
		profileLastName:        "",
		profilePortalAdminOrgs: []string{},
	}

	user, _, err := c.c.User.CreateUser(ctx, okta.CreateUserRequest{Profile: &profile}, &query.Params{Activate: &activate})
	if err != nil {
		var cErr *okta.Error
		if errors.As(err, &cErr) {
			return nil, convertOktaError(cErr)
		}

		return nil, errors.Wrapf(err, "failed to register user")
	}

	return convertUser(user)
}

// UpdateUser updates the Okta user. It takes UpdateProfileParams and apply them to the user with the given userID.
// Returns the updated User and an error.
func (c *Client) UpdateUser(ctx context.Context, userID string, params UpdateUserParams) (*User, error) {
	l := extractLogger(ctx)

	zapVal := zap.Skip()
	if l.Core().Enabled(zap.DebugLevel) {
		zapVal = zap.Any("params", params)
	}

	l.Info("Validating params for Okta user profile update", zap.String("userID", userID), zapVal)
	err := validateUpdateUserParams(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate params")
	}

	l.Info("Updating Okta user profile.", zap.String("userID", userID))
	userToUpdate, _, err := c.c.User.GetUser(ctx, url.QueryEscape(userID))
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrap(err, "failed to find user")
	}

	// apply params
	newProfile := updatedProfile(*userToUpdate.Profile, params)
	userToUpdate.Profile = &newProfile

	// update user
	updatedUser, _, err := c.c.User.UpdateUser(ctx, userID, *userToUpdate, nil)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrap(err, "failed to find user")
	}

	return convertUser(updatedUser)
}

// SignIn returns user id and session oktaAPIToken. Returns AuthError in case of invalid login or password.
func (c *Client) SignIn(ctx context.Context, login, password string) (string, string, error) {
	l := extractLogger(ctx)
	l.Info("Signin in Okta user.", zap.String("username", login))

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

// SuspendUser set user status in Okta to SUSPENDED.
func (c *Client) SuspendUser(ctx context.Context, userID string) error {
	_, err := c.c.User.SuspendUser(ctx, userID)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to suspend user")
	}

	return nil
}

// DeleteUser deactivates and deletes user.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	l := extractLogger(ctx)
	l.Info("Deleting Okta user by ID.", zap.String("userID", userID))

	err := c.deactivateUser(ctx, userID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		l.Error("Okta user is not found.", zap.Error(err))
		return err
	}

	if err == nil {
		nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := c.waitForDeactivation(nCtx, userID); err != nil {
			l.Error("User deactivation failed.", zap.Error(err))
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

// GetRegisteredUsersCount returns number of registered users.
func (c *Client) GetRegisteredUsersCount(ctx context.Context) (float64, error) {
	l := extractLogger(ctx)
	l.Info("Getting registered users count from Okta.")

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
	l := extractLogger(ctx)
	l.Info("Creating Okta user session.")
	session, _, err := c.c.Session.CreateSession(ctx, okta.CreateSessionRequest{SessionToken: sessionToken})
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to create session")
	}

	return session.Id, *session.ExpiresAt, nil
}

// CheckSession returns user login if session exists and active.
func (c *Client) CheckSession(ctx context.Context, sessionID string) (string, error) {
	l := extractLogger(ctx)
	l.Info("Checking Okta user session.")
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
	l := extractLogger(ctx)
	l.Info("Refreshing Okta user session.")

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
	l := extractLogger(ctx)
	l.Info("Closing Okta user session.")

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
	l := extractLogger(ctx)
	l.Info("Creating Okta group.", zap.String("oktaGroupName", name))

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
	l := extractLogger(ctx)
	l.Info("Looking for Okta group by name.", zap.String("oktaGroupName", name))

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
	l := extractLogger(ctx)
	l.Info("Deleting Okta group by ID.", zap.String("oktaGroupID", groupID))

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
	l := extractLogger(ctx)
	l.Info("Looking for Okta group members.", zap.String("groupID", groupID))

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
		resultUser, err := convertUser(user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert user type")
		}

		res = append(res, *resultUser)
	}

	return res, nil
}

// AddUserToGroup add user to group.
func (c *Client) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	l := extractLogger(ctx)
	l.Info("Adding user to Okta group.",
		zap.String("oktaUserID", userID),
		zap.String("groupID", groupID),
	)

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
var ErrOriginNotFound error = errors.New("trusted origin was not found") //nolint:revive

// GetTrustedOriginID returns origin's id if it exists, nil and error when it does not.
func (c *Client) GetTrustedOriginID(ctx context.Context, origin string) (string, error) {
	var origins []*okta.TrustedOrigin
	err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/trustedOrigins?q=%s", url.QueryEscape(origin)), nil, &origins)
	if err != nil {
		return "", errors.Wrap(err, "failed to get origin")
	}

	for _, trusted := range origins {
		if trusted.Origin == origin {
			return trusted.Id, nil
		}
	}
	return "", ErrOriginNotFound
}

// CreateTrustedOrigin makes the given origin trusted and returns it's id.
func (c *Client) CreateTrustedOrigin(ctx context.Context, origin string) (string, error) {
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
	l := extractLogger(ctx)
	l.Info("Deleting user from Okta group.",
		zap.String("oktaUserID", userID),
		zap.String("groupID", groupID),
	)
	_, err := c.c.Group.RemoveUserFromGroup(ctx, groupID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user to group")
	}

	return nil
}

// ResetPassword resets user password and sends email for setting new one.
func (c *Client) ResetPassword(ctx context.Context, userID string) error {
	l := extractLogger(ctx)
	l.Info("Resetting password for Okta user.", zap.String("oktaUserID", userID))

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
	l := extractLogger(ctx)
	l.Info("Updating Okta user profile schema.",
		zap.String("typeID", typeID),
		zap.String("schemaID", schema.ID),
	)

	var res Schema
	err := c.DoRequest(ctx, "POST", "/api/v1/meta/schemas/user/"+typeID, schema, &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update schema")
	}

	return &res, nil
}

// GetSchema returns schema for provided type.
func (c *Client) GetSchema(ctx context.Context, typeID string) (*Schema, error) {
	l := extractLogger(ctx)
	l.Info("Looking for Okta user profile schema.", zap.String("typeID", typeID))
	var schema Schema
	err := c.DoRequest(ctx, "GET", "/api/v1/meta/schemas/user/"+typeID, nil, &schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schema")
	}

	return &schema, nil
}

func (c *Client) ListPolicies(ctx context.Context, qp *query.Params) ([]*okta.Policy, error) { //nolint:revive
	l := extractLogger(ctx)
	l.Info("Looking for Okta policies.")
	policyInterfaces, _, err := c.c.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return nil, err
	}

	policies := make([]*okta.Policy, 0, len(policyInterfaces))
	for _, policyInterface := range policyInterfaces {
		if !policyInterface.IsPolicyInstance() {
			continue
		}

		policy, ok := policyInterface.(*okta.Policy)
		if !ok {
			continue
		}

		policies = append(policies, policy)
	}

	return policies, err
}

const createOAuthAppRequestBody = `
{
    "name": "oidc_client",
    "label": "PMM-Grafana-{{ .PMMServerID }}",
    "status": "ACTIVE",
    "signOnMode": "OPENID_CONNECT",
    "credentials": {
        "oauthClient": {
            "token_endpoint_auth_method": "none"
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
                "authorization_code",
                "refresh_token"
            ],
            "application_type": "browser",
            "consent_method": "REQUIRED"
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

const createMachineAuthAppRequestBody = `
{
    "name": "oidc_client",
    "label": "PMM-Managed-{{ .PMMServerID }}",
    "signOnMode": "OPENID_CONNECT",
    "credentials": {
        "oauthClient": {
            "token_endpoint_auth_method": "client_secret_basic"
        }
    },
    "settings": {
        "oauthClient": {
            "application_type": "service",
            "consent_method": "REQUIRED",
            "grant_types": [
                "client_credentials"
            ],
            "response_types": [
                "token"
            ]
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

// OAuthApp represents an oauth app.
type OAuthApp struct {
	AppID       string `json:"id"`
	Credentials struct {
		OAuthClient struct {
			ClientID string `json:"client_id"` // nolint:tagliatelle
		} `json:"oauthClient"`
	} `json:"credentials"`
}

// MachineAuthApp represents a machine-to-machine authorized app.
type MachineAuthApp struct {
	AppID       string `json:"id"`
	Credentials struct {
		OAuthClient struct {
			ClientID     string `json:"client_id"`     // nolint:tagliatelle
			ClientSecret string `json:"client_secret"` // nolint:tagliatelle
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

// MachineAuthAppParams contains values needed when creating a new machine-to-machine app.
type MachineAuthAppParams struct {
	PMMServerID string
	OrgID       string
	InventoryID string
}

// CreateOAuthApp creates a new OAuth app.
func (c *Client) CreateOAuthApp(ctx context.Context, params *OAuthAppParams) (*OAuthApp, error) {
	var request bytes.Buffer
	tmpl := template.Must(template.New("CreateOAuthAppRequest").Parse(createOAuthAppRequestBody))
	err := tmpl.Execute(&request, params)
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

// CreateMachineAuthApp creates a new OAuth app.
func (c *Client) CreateMachineAuthApp(ctx context.Context, params *MachineAuthAppParams) (*MachineAuthApp, error) {
	var request bytes.Buffer
	tmpl := template.Must(template.New("CreateMachineAuthAppRequest").Parse(createMachineAuthAppRequestBody))
	err := tmpl.Execute(&request, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request body for adding the OAuth App")
	}

	var result MachineAuthApp
	err = c.DoRequest(ctx, "POST", "/api/v1/apps", &request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "adding OAuth APP to Okta failed")
	}
	return &result, nil
}

// DeleteApp deletes an app with given appID.
func (c *Client) DeleteApp(ctx context.Context, appID string) error {
	l := logger.GetLoggerFromContext(ctx).Named("oktaClient")
	_, err := c.c.Application.DeactivateApplication(ctx, appID)
	if err != nil {
		return errors.Wrap(err, "failed to deactivate app before deleting it")
	}

	if _, err = c.c.Application.DeleteApplication(ctx, appID); err != nil {
		_, e := c.c.Application.ActivateApplication(ctx, appID)
		if e != nil {
			l.Error("Failed to re-activate app after deleting failed, manual intervention in Okta required", zap.Error(e))
		}
		return err
	}
	return nil
}

// GetActivationLink returns activation url for users that are not activated yet.
func (c *Client) GetActivationLink(ctx context.Context, userID string) (string, error) {
	l := logger.GetLoggerFromContext(ctx).Named("oktaClient")
	sendEmail := false
	activationInfo, _, err := c.c.User.ActivateUser(ctx, userID, &query.Params{SendEmail: &sendEmail})
	if err != nil {
		l.Error("Failed to get activation link", zap.Error(err))
		return "", errors.Wrap(err, "failed to activate user")
	}

	return activationInfo.ActivationUrl, nil
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
	login, ok := profile[profileLogin]
	if !ok {
		return "", errors.New("missing user " + profileLogin)
	}

	result, ok := login.(string)
	if !ok {
		result = ""
	}

	return result, nil
}

func getUserFirstName(user *okta.User) (string, error) {
	if user.Profile == nil {
		return "", errors.New("missing user profile")
	}

	profile := *user.Profile
	name, ok := profile[profileFirstName]
	if !ok {
		return "", errors.New("missing user " + profileFirstName)
	}

	result, ok := name.(string)
	if !ok {
		result = ""
	}

	return result, nil
}

func getUserLastName(user *okta.User) (string, error) {
	if user.Profile == nil {
		return "", errors.New("missing user profile")
	}

	profile := *user.Profile
	name, ok := profile[profileLastName]
	if !ok {
		return "", errors.New("missing user " + profileLastName)
	}

	result, ok := name.(string)
	if !ok {
		result = ""
	}

	return result, nil
}

func getPortalAdminOrgs(user *okta.User) ([]string, error) {
	if user.Profile == nil {
		return nil, errors.New("missing user profile")
	}

	result := []string{}

	profile := *user.Profile
	orgsRaw, ok := profile[profilePortalAdminOrgs]
	if !ok {
		return result, nil
	}

	orgsSlice, ok := orgsRaw.([]interface{})
	if !ok {
		return result, nil
	}

	for _, val := range orgsSlice {
		if str, ok := val.(string); ok {
			result = append(result, str)
		}
	}

	return result, nil
}

func extractLogger(ctx context.Context) *zap.Logger {
	return logger.GetLoggerFromContext(ctx).Named("oktaClient")
}

func convertUser(oktaUser *okta.User) (*User, error) {
	userLogin, err := getUserLogin(oktaUser)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", oktaUser.Id)
	}

	firstName, err := getUserFirstName(oktaUser)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", oktaUser.Id)
	}

	lastName, err := getUserLastName(oktaUser)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", oktaUser.Id)
	}

	portalAdminOrgs, err := getPortalAdminOrgs(oktaUser)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s has bad profile", oktaUser.Id)
	}

	return &User{
		ID:              oktaUser.Id,
		Login:           userLogin,
		FirstName:       firstName,
		LastName:        lastName,
		Status:          oktaUser.Status,
		PortalAdminOrgs: portalAdminOrgs,
	}, nil
}

func updatedProfile(profile okta.UserProfile, params UpdateUserParams) okta.UserProfile {
	if params.PortalAdminOrgs != nil {
		profile[profilePortalAdminOrgs] = params.PortalAdminOrgs
	}

	if params.Firstname != nil {
		profile[profileFirstName] = params.Firstname
	}

	if params.Lastname != nil {
		profile[profileLastName] = params.Lastname
	}

	return profile
}

func validateUpdateUserParams(params UpdateUserParams) error {
	if params.PortalAdminOrgs != nil {
		ids := *params.PortalAdminOrgs

		err := validatePortalAdminOrgs(ids)
		if err != nil {
			return err
		}
	}

	if params.Firstname != nil && *params.Firstname == "" {
		return ErrEmptyFirstName
	}

	if params.Lastname != nil && *params.Lastname == "" {
		return ErrEmptyLastName
	}

	return nil
}

func validatePortalAdminOrgs(ids []string) error {
	// map to check duplicates
	duplMap := make(map[string]struct{}, len(ids))
	var err error

	for _, val := range ids {
		_, err = uuid.Parse(val)
		if err != nil {
			return ErrInvalidPortalAdminOrgs
		}

		if _, ok := duplMap[val]; ok {
			return ErrDuplicatedPortalAdminOrgs
		}

		duplMap[val] = struct{}{}
	}

	return nil
}

func validateRegisterUserParams(params RegisterUserParams) error {
	if params.Login == "" {
		return ErrEmptyLogin
	}

	return nil
}
