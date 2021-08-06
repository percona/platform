// Package okta implements methods for interacting with Okta API.
package okta

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/logger"
)

const (
	sessionTTL         = 7 * 24 * time.Hour
	sessionTTLRuleName = "SessionTTL (Managed by Auth Service. DO NOT EDIT!)"
)

// Service implements methods for interacting with Okta API. Some methods can return
// // AuthError which indicates on authentication/authorisation problems.
type Service struct {
	l     *zap.SugaredLogger
	c     *okta.Client
	host  string
	token string
}

// New returns new Service instance.
func New(host, token string) (*Service, error) {
	l := zap.L().Named("okta").Sugar()

	httpClient := http.Client{
		Transport: logger.HTTP(http.DefaultTransport, l.Debugf),
	}

	u := url.URL{Scheme: "https", Host: host}

	_, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(u.String()),
		okta.WithToken(token),
		okta.WithHttpClient(httpClient),
		okta.WithCache(false))
	if err != nil {
		return nil, err
	}

	return &Service{
		l:     l,
		c:     client,
		host:  host,
		token: token,
	}, nil
}

// SignUp creates user account with provided login, first name and last name.
// The user is created with a "PROVISIONED" state and a verification email is sent to the user.
// Once the user sets the password the state changes to "ACTIVE" in Okta.
// Returns AuthError when login, firstName and lastName violates validation rules, also when login already exists.
func (s *Service) SignUp(ctx context.Context, login, firstName, lastName string) (*User, error) {
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
	user, _, err := s.c.User.CreateUser(ctx, u, qp)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
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
func (s *Service) FindUser(ctx context.Context, login string) (*okta.User, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	user, _, err := s.c.User.GetUser(ctx, url.QueryEscape(login))
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrapf(err, "failed to find user")
	}

	return user, nil
}

// SignIn returns user id and session token. Returns AuthError in case of invalid login or password.
func (s *Service) SignIn(ctx context.Context, login, password string) (string, string, error) {
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
		} `json:"_embedded"`
	}{}

	err := s.doRequest(ctx, "POST", "/api/v1/authn", data, &resp)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return "", "", convertOktaError(oErr)
		}

		return "", "", errors.Wrap(err, "failed to authenticate user")
	}

	return resp.Embedded.User.ID, resp.SessionToken, nil
}

// DeleteUser deactivates and deletes user.
func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	err := s.deactivateUser(ctx, userID)
	if err != nil && err != ErrNotFound {
		return err
	}

	if err == nil {
		nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := s.waitForDeactivation(nCtx, userID); err != nil {
			return err
		}
	}

	_, err = s.c.User.DeactivateOrDeleteUser(ctx, userID, nil)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}

// UpdateProfile updates user's profile.
func (s *Service) UpdateProfile(ctx context.Context, user *okta.User, firstName, lastName string) (*okta.User, error) {
	login, err := getUserLogin(user)
	if err != nil {
		return nil, err
	}

	// The okta-go-sdk implementation differs slightly from the API docs where it uses
	// PUT instead of POST for this endpoint. In case of PUT (used by the SDK) it requires
	// both the email and login fields to be non-empty, so we send the original login to prevent
	// the request from returning an error.
	body := okta.User{ //nolint:exhaustivestruct
		Profile: &okta.UserProfile{
			"firstName": firstName,
			"lastName":  lastName,
			"login":     login,
			"email":     login,
		},
	}
	u, _, err := s.c.User.UpdateUser(ctx, user.Id, body, nil)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return nil, convertOktaError(oErr)
		}

		return nil, errors.Wrapf(err, "failed to update user profile")
	}

	return u, nil
}

func (s *Service) deactivateUser(ctx context.Context, userID string) error {
	_, err := s.c.User.DeactivateUser(ctx, userID, nil)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to deactivate user")
	}

	return nil
}

func (s *Service) waitForDeactivation(ctx context.Context, userID string) error {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		user, _, err := s.c.User.GetUser(ctx, userID)
		if err != nil {
			if oErr, ok := err.(*okta.Error); ok {
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
func (s *Service) GetRegisteredUsersCount(ctx context.Context) (float64, error) {
	qp := query.NewQueryParams(query.WithQ("Everyone"), query.WithFilter("type eq \"BUILT_IN\""))
	groups, _, err := s.c.Group.ListGroups(ctx, qp)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return 0, convertOktaError(oErr)
		}

		return 0, errors.Wrap(err, "failed to find everyone users group")
	}

	if len(groups) != 1 {
		return 0, fmt.Errorf("expect only one everyone group search result, got %d", len(groups))
	}

	var group okta.Group
	path := fmt.Sprintf("/api/v1/groups/%s?expand=stats", groups[0].Id)
	err = s.doRequest(ctx, "GET", path, nil, &group)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
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
func (s *Service) CreateSession(ctx context.Context, sessionToken string) (string, time.Time, error) {
	session, _, err := s.c.Session.CreateSession(ctx, okta.CreateSessionRequest{SessionToken: sessionToken})
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "failed to create session")
	}

	return session.Id, *session.ExpiresAt, nil
}

// CheckSession returns user login if session exists and active.
func (s *Service) CheckSession(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", ErrNotFound
	}

	session, _, err := s.c.Session.GetSession(ctx, sessionID)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return "", convertOktaError(oErr)
		}

		return "", errors.Wrap(err, "failed to check session")
	}

	return session.Login, nil
}

// RefreshSession resets session timeout and returns new expiration date.
func (s *Service) RefreshSession(ctx context.Context, sessionID string) (time.Time, error) {
	session, _, err := s.c.Session.RefreshSession(ctx, sessionID)
	if err != nil || session.ExpiresAt == nil {
		if oErr, ok := err.(*okta.Error); ok {
			return time.Time{}, convertOktaError(oErr)
		}

		return time.Time{}, errors.Wrap(err, "failed to refresh session")
	}

	return *session.ExpiresAt, nil
}

// CloseSession terminates user session.
func (s *Service) CloseSession(ctx context.Context, sessionID string) error {
	_, err := s.c.Session.EndSession(ctx, sessionID)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to close session")
	}

	return nil
}

// CreateGroup creates group with provided name and description.
func (s *Service) CreateGroup(ctx context.Context, name, description string) (*Group, error) {
	req := okta.Group{
		Profile: &okta.GroupProfile{
			Name:        name,
			Description: description,
		},
	}

	group, _, err := s.c.Group.CreateGroup(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create group")
	}

	return &Group{
		ID:          group.Id,
		Name:        group.Profile.Name,
		Description: group.Profile.Description,
	}, nil
}

// DeleteGroup delete group with provided ID.
func (s *Service) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := s.c.Group.DeleteGroup(ctx, groupID)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to delete group")
	}

	return nil
}

// GetGroupMembers returns list of group members.
func (s *Service) GetGroupMembers(ctx context.Context, groupID string, limit int, cursor string) ([]User, error) {
	params := query.Params{
		Limit:  int64(limit),
		Cursor: cursor,
	}
	users, _, err := s.c.Group.ListGroupUsers(ctx, groupID, &params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get group members")
	}

	res := make([]User, 0, len(users))
	for _, user := range users {
		login, err := getUserLogin(user)
		if err != nil {
			s.l.Warnf("User %s has bad profile, reason: %+v.", user.Id, err)
		}

		res = append(res, User{ID: user.Id, Login: login, Status: user.Status})
	}

	return res, nil
}

// AddUserToGroup add user to group.
func (s *Service) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	_, err := s.c.Group.AddUserToGroup(ctx, groupID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user to group")
	}

	return nil
}

// RemoveUserToGroup remove user from group.
func (s *Service) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	_, err := s.c.Group.RemoveUserFromGroup(ctx, groupID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user to group")
	}

	return nil
}

// ResetPassword resets user password and sends email for setting new one.
func (s *Service) ResetPassword(ctx context.Context, userID string) error {
	qp := query.NewQueryParams(query.WithSendEmail(true))
	_, _, err := s.c.User.ResetPassword(ctx, userID, qp)
	if err != nil {
		if oErr, ok := err.(*okta.Error); ok {
			return convertOktaError(oErr)
		}

		return errors.Wrap(err, "failed to reset password")
	}

	return nil
}

// UpdateSchema updates schema for provided type.
func (s *Service) UpdateSchema(ctx context.Context, typeID string, schema *Schema) (*Schema, error) {
	var res Schema
	err := s.doRequest(ctx, "POST", "/api/v1/meta/schemas/user/"+typeID, schema, &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update schema")
	}

	return &res, nil
}

// GetSchema returns schema for provided type.
func (s *Service) GetSchema(ctx context.Context, typeID string) (*Schema, error) {
	var schema Schema
	err := s.doRequest(ctx, "GET", "/api/v1/meta/schemas/user/"+typeID, nil, &schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schema")
	}

	return &schema, nil
}

// MakeFirstAndLastNamesOptional makes firsName and LastName fields optional in provided schema.
func (s *Service) MakeFirstAndLastNamesOptional(ctx context.Context, typeID string) error {
	schema := Schema{
		Definitions: map[string]Definition{
			"base": {
				ID:   "#base",
				Type: "object",
				Properties: map[string]DefinitionProperty{
					"firstName": {
						Required: pointer.ToBool(false),
					},
					"lastName": {
						Required: pointer.ToBool(false),
					},
				},
			},
		},
	}

	_, err := s.UpdateSchema(ctx, typeID, &schema)
	if err != nil {
		return errors.Wrap(err, "failed to make lastName and firstName fields optional")
	}

	return nil
}

// SetUpSessionTTLRule checks whether session TTL rule exists. If it's not then rule created,
// if rule exists but TTL value is not actual then it updated.
func (s *Service) SetUpSessionTTLRule(ctx context.Context) error {
	policyID, err := s.findDefaultPolicy(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get default policy")
	}

	// Get policy rules
	var rules []*signOnPolicyRule
	if err := s.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/policies/%s/rules", policyID), nil, &rules); err != nil {
		return errors.Wrap(err, "failed to get default policy rules")
	}

	for _, rule := range rules {
		if rule.Name == sessionTTLRuleName {
			// Check session TTL value
			if rule.Actions.SignOn.Session.MaxSessionIdleMinutes == int64(sessionTTL.Minutes()) {
				return nil
			}

			// Update session TTL value
			rule.Actions.SignOn.Session.MaxSessionIdleMinutes = int64(sessionTTL.Minutes())
			if err := s.doRequest(ctx, "PUT", fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, rule.ID), rule, nil); err != nil {
				return errors.Wrap(err, "failed to update session TTL rule")
			}

			return nil
		}
	}

	// Create session TTL rule
	sessionTTLRule := signOnPolicyRule{
		Name: sessionTTLRuleName,
		Type: "SIGN_ON",
		Conditions: &okta.OktaSignOnPolicyRuleConditions{
			AuthContext: &okta.PolicyRuleAuthContextCondition{
				AuthType: "ANY",
			},
			Network: &okta.PolicyNetworkCondition{
				Connection: "ANYWHERE",
			},
		},
		Actions: &signOnPolicyRuleActions{
			SignOn: &signOnPolicyRuleSignOnActions{
				Access: "ALLOW",
				Session: &signOnPolicyRuleSignOnSessionActions{
					MaxSessionIdleMinutes: int64(sessionTTL.Minutes()),
					UsePersistentCookie:   pointer.ToBool(false),
				},
			},
		},
	}

	if err := s.doRequest(ctx, "POST", fmt.Sprintf("/api/v1/policies/%s/rules", policyID), sessionTTLRule, nil); err != nil {
		return errors.Wrap(err, "failed to create session TTL rule")
	}

	return nil
}

// findDefaultPolicy returns default Okta Sign On policy ID.
func (s *Service) findDefaultPolicy(ctx context.Context) (string, error) {
	qp := query.NewQueryParams(query.WithType("OKTA_SIGN_ON"))
	policies, _, err := s.c.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return "", err
	}

	for _, policy := range policies {
		if *policy.System && policy.Name == "Default Policy" {
			return policy.Id, nil
		}
	}

	return "", errors.New("can't find default Okta policy")
}

func (s *Service) doRequest(ctx context.Context, method, path string, body, v interface{}) error {
	requestExecutor := s.c.GetRequestExecutor().WithAccept("application/json").WithContentType("application/json")
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

func convertOktaError(error *okta.Error) error {
	switch error.ErrorCode {
	case "E0000001":
		switch error.ErrorSummary {
		case "Api validation failed: password":
			return &AuthError{msg: "invalid password", origin: error}
		case "Api validation failed: login":
			return &AuthError{msg: "invalid login", origin: error}
		default:
			return error
		}
	case "E0000004":
		return ErrAuthentication
	case "E0000007":
		return ErrNotFound
	default:
		return error
	}
}
