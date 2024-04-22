package okta

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/okta/okta-sdk-golang/v2/okta"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var authErrorType = new(AuthError) //nolint:gochecknoglobals,errname

const sessionTTL = 7 * 24 * time.Hour

func init() { //nolint:gochecknoinits
	gofakeit.Seed(time.Now().UnixNano())
}

func createOktaService(t *testing.T) (*Client, error) {
	t.Helper()
	return New(context.Background(), OktaDevHost, GetOktaToken(t))
}

func TestSignUp(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("invalid login", func(t *testing.T) {
		t.Parallel()

		_, _, firstName, lastName := GenCredentials(t)
		user, err := s.SignUp(context.Background(), "not email", firstName, lastName)
		require.EqualError(t, err, "invalid login: Api validation failed: login")
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("empty login", func(t *testing.T) {
		t.Parallel()

		_, _, firstName, lastName := GenCredentials(t)
		user, err := s.SignUp(context.Background(), "", firstName, lastName)
		require.Equal(t, err, ErrEmptyLogin)
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("empty first name", func(t *testing.T) {
		t.Parallel()

		email, _, _, lastName := GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, "", lastName)
		require.Equal(t, err, ErrEmptyFirstName)
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("empty last name", func(t *testing.T) {
		t.Parallel()

		email, _, firstName, _ := GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, "")
		require.Equal(t, err, ErrEmptyLastName)
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("valid sign up", func(t *testing.T) {
		t.Parallel()

		email, _, firstName, lastName := GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, lastName)
		require.NoError(t, err)
		defer DeleteUser(t, user.ID)

		require.Equal(t, email, user.Login)
		require.Equal(t, UserStatusProvisioned, user.Status)
		require.NotEmpty(t, user.ID)
	})
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	t.Run("invalid password", func(t *testing.T) {
		t.Parallel()

		userID, sessionToken, err := s.SignIn(context.Background(), email, "wrong")
		require.Equal(t, ErrAuthentication, err)
		require.IsType(t, authErrorType, err)
		require.Empty(t, sessionToken)
		require.Empty(t, userID)
	})

	t.Run("empty password", func(t *testing.T) {
		t.Parallel()

		userID, sessionToken, err := s.SignIn(context.Background(), email, "")
		require.Equal(t, ErrEmptyPassword, err)
		require.IsType(t, authErrorType, err)
		require.Empty(t, sessionToken)
		require.Empty(t, userID)
	})

	t.Run("invalid login", func(t *testing.T) {
		t.Parallel()

		userID, sessionToken, err := s.SignIn(context.Background(), "wrong", password)
		require.Equal(t, ErrAuthentication, err)
		require.IsType(t, authErrorType, err)
		require.Empty(t, sessionToken)
		require.Empty(t, userID)
	})

	t.Run("empty login", func(t *testing.T) {
		t.Parallel()

		userID, sessionToken, err := s.SignIn(context.Background(), "", password)
		require.Equal(t, err, ErrEmptyLogin)
		require.IsType(t, authErrorType, err)
		require.Empty(t, sessionToken)
		require.Empty(t, userID)
	})

	t.Run("valid sign in", func(t *testing.T) {
		t.Parallel()

		userID, sessionToken, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)
		require.NotEmpty(t, sessionToken)
		require.NotEmpty(t, userID)
	})
}

func TestSignInByToken(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		t.Parallel()

		email, password, firstName, lastName := GenCredentials(t)
		user := CreateInactivatedTestUser(t, email, password, firstName, lastName)

		t.Cleanup(func() {
			DeleteUser(t, user.ID)
		})

		token := ActivateUser(t, user.ID)

		authInfo, err := s.SignInByToken(context.Background(), token)
		require.NoError(t, err)
		require.NotEmpty(t, authInfo)
		require.Equal(t, user.ID, authInfo.Embedded.User.ID)
	})

	t.Run("wrong token", func(t *testing.T) {
		t.Parallel()

		authInfo, err := s.SignInByToken(context.Background(), gofakeit.UUID())
		require.Error(t, err)
		require.ErrorContains(t, err, "authentication error")
		require.Empty(t, authInfo)
	})
}

func TestSignInByStateToken(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		t.Parallel()

		email, _, firstName, lastName := GenCredentials(t)
		user := CreateInactivatedTestUser(t, email, "", firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, user.ID)
		})

		activationToken := ActivateUser(t, user.ID)

		data := struct {
			Token string `json:"token"`
		}{Token: activationToken}

		resp := AuthenticatedInfo{}
		client := createOktaClient(t)
		err = oktaAPIRequest(client, "POST", "/api/v1/authn", data, &resp)
		require.NoError(t, err)
		require.NotEmpty(t, resp.StateToken)

		authInfo, err := s.SignInByStateToken(context.Background(), resp.StateToken)
		require.NoError(t, err)
		require.NotEmpty(t, authInfo)
		require.Equal(t, user.ID, authInfo.Embedded.User.ID)
	})

	t.Run("wrong token", func(t *testing.T) {
		t.Parallel()

		authInfo, err := s.SignInByStateToken(context.Background(), gofakeit.UUID())
		require.Error(t, err)
		require.ErrorContains(t, err, "Invalid token provided")
		require.Empty(t, authInfo)
	})
}

func TestSessions(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	t.Run("invalid session", func(t *testing.T) {
		t.Parallel()

		login, err := s.CheckSession(context.Background(), "invalid-session-oktaAPIToken")
		require.Equal(t, err, ErrNotFound)
		require.IsType(t, authErrorType, err)
		require.Empty(t, login)
	})

	t.Run("valid session", func(t *testing.T) {
		t.Parallel()

		userID, token, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, userID)

		ts := time.Now()
		timeError := 5 * time.Second
		sessionID, expiresAt, err := s.CreateSession(context.Background(), token)
		require.NoError(t, err)
		require.NotEmpty(t, sessionID)
		require.NotEmpty(t, expiresAt)
		require.GreaterOrEqual(t, expiresAt.Unix(), ts.Add(sessionTTL-timeError).Unix())
		require.LessOrEqual(t, expiresAt.Unix(), time.Now().Add(sessionTTL+timeError).Unix())

		userEmail, err := s.CheckSession(context.Background(), sessionID)
		require.NoError(t, err)
		require.Equal(t, email, userEmail)
	})
}

func TestSessionRefresh(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		_, token, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)

		sessionID, expiresAt, err := s.CreateSession(context.Background(), token)
		require.NoError(t, err)
		require.NotEmpty(t, expiresAt)

		time.Sleep(time.Second)

		newExpirationTime, err := s.RefreshSession(context.Background(), sessionID)
		require.NoError(t, err)
		require.NotEmpty(t, newExpirationTime)

		//nolint:godox
		// TODO: https://jira.percona.com/browse/PMM-12756
		// Requires investigation, it seems that Okta changed the behavior of refresh tokens:
		// https://developer.okta.com/docs/guides/refresh-tokens/main/#enable-refresh-token-rotation
		// > Note: When a refresh token is rotated, the new refresh_token string in the response has a different value than the previous refresh_token string due to security concerns with single-page apps. However, the expiration date remains the same. The lifetime is inherited from the initial refresh token minted when the user first authenticates.
		// require.Greater(t, newExpirationTime.Unix(), expiresAt.Unix())
	})

	t.Run("invalid session", func(t *testing.T) {
		t.Parallel()

		expTime, err := s.RefreshSession(context.Background(), "invalid-session-id")
		require.Equal(t, err, ErrNotFound)
		require.Zero(t, expTime)
	})
}

func TestCloseSession(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		_, token, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)

		sessionID, _, err := s.CreateSession(context.Background(), token)
		require.NoError(t, err)

		err = s.CloseSession(context.Background(), sessionID)
		require.NoError(t, err)

		_, err = s.CheckSession(context.Background(), sessionID)
		require.Equal(t, err, ErrNotFound)
		require.IsType(t, authErrorType, err)
	})

	t.Run("invalid session", func(t *testing.T) {
		t.Parallel()

		err = s.CloseSession(context.Background(), "invalid-session-id")
		require.Equal(t, err, ErrNotFound)
	})

	t.Run("already closed session", func(t *testing.T) {
		t.Parallel()

		_, token, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)

		sessionID, _, err := s.CreateSession(context.Background(), token)
		require.NoError(t, err)

		err = s.CloseSession(context.Background(), sessionID)
		require.NoError(t, err)

		err = s.CloseSession(context.Background(), sessionID)
		require.Equal(t, err, ErrNotFound)
	})
}

func TestFindUser(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	t.Run("user doesn't exists", func(t *testing.T) {
		t.Parallel()

		userID, err := s.FindUser(context.Background(), "invalid@example.com")
		require.Equal(t, ErrNotFound, err)
		require.Empty(t, userID)
	})

	t.Run("user exists", func(t *testing.T) {
		t.Parallel()

		foundUser, err := s.FindUser(context.Background(), email)
		require.NoError(t, err)
		require.NotEmpty(t, foundUser)
		require.Equal(t, []string{}, foundUser.PortalAdminOrgs)
		require.Equal(t, firstName, foundUser.FirstName)
		require.Equal(t, lastName, foundUser.LastName)
	})
}

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		email, _, _, _ := GenCredentials(t)

		ctx := context.Background()

		u, err := s.RegisterUser(ctx, RegisterUserParams{Login: email})
		require.NoError(t, err)
		t.Cleanup(func() {
			DeleteUser(t, u.ID)
		})

		user, _, err := s.c.User.GetUser(ctx, u.ID)
		require.NoError(t, err)
		require.Equal(t, UserStatusProvisioned, user.Status)

		require.Equal(t, u.Login, email)
	})

	t.Run("invalid email", func(t *testing.T) {
		t.Parallel()

		_, err := s.RegisterUser(context.Background(), RegisterUserParams{Login: "not_an_email"})
		require.EqualError(t, err, "invalid login: Api validation failed: login")
	})

	t.Run("user exists", func(t *testing.T) {
		t.Parallel()

		email, password, firstName, lastName := GenCredentials(t)
		user := CreateTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, user.ID)
		})

		_, err := s.RegisterUser(context.Background(), RegisterUserParams{Login: email})
		require.EqualError(t, err, "invalid login: Api validation failed: login")
	})
}

func TestRegisterInactiveUser(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		email, _, _, _ := GenCredentials(t)
		ctx := context.Background()

		u, err := s.RegisterInactiveUser(ctx, RegisterUserParams{Login: email})
		require.NoError(t, err)
		t.Cleanup(func() {
			DeleteUser(t, u.ID)
		})

		user, _, err := s.c.User.GetUser(ctx, u.ID)
		require.NoError(t, err)
		require.Equal(t, UserStatusStaged, user.Status)

		require.Equal(t, u.Login, email)
	})
}

func TestPasswordReset(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	u, err := s.FindUser(context.Background(), email)
	require.NoError(t, err)

	err = s.ResetPassword(context.Background(), u.ID)
	require.NoError(t, err)

	_, _, err = s.SignIn(context.Background(), email, password)
	require.Equal(t, ErrAuthentication, err)
}

func TestGroups(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	user := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, user.ID)
	})

	name := gofakeit.LastName() + ", " + gofakeit.LastName() + " and " + gofakeit.LastName()
	description := "Test group"
	group, err := s.CreateGroup(context.Background(), name, description)
	t.Cleanup(func() {
		DeleteGroup(t, group.ID)
	})
	require.NoError(t, err)
	require.Equal(t, name, group.Name)
	require.Equal(t, description, group.Description)
	require.NotEmpty(t, group.ID)

	err = s.AddUserToGroup(context.Background(), user.ID, group.ID)
	require.NoError(t, err)

	users, err := s.GetGroupMembers(context.Background(), group.ID, 0, "")
	require.NoError(t, err)

	require.Len(t, users, 1)
	require.Equal(t, *user, users[0])

	exists, err := s.GroupExists(context.Background(), name)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = s.GroupExists(context.Background(), "non-existent-group")
	require.NoError(t, err)
	require.False(t, exists)
}

func TestFindGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s, err := createOktaService(t)
	require.NoError(t, err)

	name := gofakeit.LastName() + ", " + gofakeit.LastName() + " and " + gofakeit.LastName()
	description := "Test group"

	groups, err := s.FindGroupByName(ctx, name)
	require.NoError(t, err)
	require.Empty(t, groups)

	group, err := s.CreateGroup(ctx, name, description)
	t.Cleanup(func() {
		DeleteGroup(t, group.ID)
	})
	require.NoError(t, err)

	groups, err = s.FindGroupByName(ctx, name)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, name, groups[0].Name)
	require.Equal(t, description, groups[0].Description)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		email, _, firstName, lastName := GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, lastName)
		require.NoError(t, err)

		err = s.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)

		_, err = s.FindUser(context.Background(), user.Login)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("missing user", func(t *testing.T) {
		t.Parallel()

		err = s.DeleteUser(context.Background(), "unknown-id")
		require.Equal(t, ErrNotFound, err)
	})
}

func TestSuspendUser(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("user suspended", func(t *testing.T) {
		t.Parallel()

		email, password, firstName, lastName := GenCredentials(t)
		user := CreateTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, user.ID)
		})

		err = s.SuspendUser(context.Background(), user.ID)
		require.NoError(t, err)

		usr, err := s.FindUser(context.Background(), user.Login)
		require.NoError(t, err)
		require.Equal(t, UserStatusSuspended, usr.Status)
	})

	t.Run("suspended user can't login", func(t *testing.T) {
		t.Parallel()

		email, password, firstName, lastName := GenCredentials(t)
		user := CreateTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, user.ID)
		})

		err = s.SuspendUser(context.Background(), user.ID)
		require.NoError(t, err)

		userID, sessionToken, err := s.SignIn(context.Background(), email, password)
		require.Equal(t, ErrAuthentication, err)
		require.IsType(t, authErrorType, err)
		require.Empty(t, sessionToken)
		require.Empty(t, userID)
	})

	t.Run("non existing user can't be suspended", func(t *testing.T) {
		t.Parallel()

		err = s.SuspendUser(context.Background(), "unknown-id")
		require.Equal(t, ErrNotFound, err)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := GenCredentials(t)
	testUser := CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		DeleteUser(t, testUser.ID)
	})

	t.Run("user doesn't exists", func(t *testing.T) {
		t.Parallel()

		otherFirstName := "firstName"
		otherLastName := "lastName"
		_, err := s.UpdateUser(context.Background(), "unknown", UpdateUserParams{Firstname: &otherFirstName, Lastname: &otherLastName})
		require.EqualError(t, err, "not found")
	})

	t.Run("user exists update lastname firstname successful", func(t *testing.T) {
		t.Parallel()

		user, err := s.FindUser(context.Background(), testUser.Login)
		require.NoError(t, err)
		t.Log(user.FirstName, user.LastName, user.Login, user.ID)

		newFirstName := gofakeit.FirstName()
		newLastName := gofakeit.LastName()
		updatedUser, err := s.UpdateUser(context.Background(), user.ID, UpdateUserParams{Firstname: &newFirstName, Lastname: &newLastName})
		require.NoError(t, err)

		require.Equal(t, user.ID, updatedUser.ID)
		require.Equal(t, newFirstName, updatedUser.FirstName)
		require.Equal(t, newLastName, updatedUser.LastName)
	})

	t.Run("user exists update portalAdminOrgs successful", func(t *testing.T) {
		t.Parallel()

		user, err := s.FindUser(context.Background(), testUser.Login)
		require.NoError(t, err)
		t.Log(user.FirstName, user.LastName, user.Login, user.ID)

		newID := uuid.NewString()
		ids := []string{newID}
		updatedUser, err := s.UpdateUser(context.Background(), user.ID, UpdateUserParams{PortalAdminOrgs: &ids})
		require.NoError(t, err)
		require.Equal(t, ids, updatedUser.PortalAdminOrgs)
		require.Equal(t, user.ID, updatedUser.ID)
	})
}

func TestGetRegisteredUsersCount(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	usersCount, err := s.GetRegisteredUsersCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, usersCount)
}

func TestAppLifecycle(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	ctx := context.Background()

	app, err := s.CreateOAuthApp(ctx, &OAuthAppParams{
		PMMServerCallbackURL: "https://localhost/graph/login/generic_oauth",
		PMMServerURL:         "https://localhost/graph",
		PMMServerID:          "0f0123ba-978d-4bcc-979d-e8495060fe81",
		OrgID:                "338311eb-3afc-45c9-b3b8-fce32f29e4e3",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		s.DeleteApp(ctx, app.AppID) //nolint:errcheck,gosec
	})
	require.NotNil(t, app)
	require.NotEmpty(t, app.AppID)
	require.NotNil(t, app.Credentials)
	require.NotNil(t, app.Credentials.OAuthClient)
	require.NotEmpty(t, app.Credentials.OAuthClient.ClientID)

	name, err := randomHex(8)
	require.NoError(t, err)

	description, err := randomHex(16)
	require.NoError(t, err)

	group, err := s.CreateGroup(ctx, name, description)
	require.NoError(t, err)
	require.NotNil(t, group)
	t.Cleanup(func() {
		s.DeleteGroup(ctx, group.ID) //nolint:errcheck,gosec
	})

	assigned := s.IsAppAssignedToGroup(ctx, app.AppID, group.ID)
	require.False(t, assigned)

	err = s.AddAppToGroup(ctx, app.AppID, group.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		s.RemoveAppFromGroup(ctx, app.AppID, group.ID) //nolint:errcheck,gosec
	})

	assigned = s.IsAppAssignedToGroup(ctx, app.AppID, group.ID)
	require.True(t, assigned)

	err = s.RemoveAppFromGroup(ctx, app.AppID, group.ID)
	require.NoError(t, err)

	assigned = s.IsAppAssignedToGroup(ctx, app.AppID, group.ID)
	require.False(t, assigned)

	err = s.DeleteGroup(ctx, group.ID)
	require.NoError(t, err)

	err = s.DeleteApp(ctx, app.AppID)
	require.NoError(t, err)
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	read, err := rand.Read(b)
	if read < n {
		return "", errors.New("failed to read given number of bytes")
	}
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func TestTrustedOrigin(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	ctx := context.Background()
	// Let's create a random domain to be able to run multiple instances of tests in parallel.
	subdomain, err := randomHex(8)
	require.NoError(t, err)
	require.NotEmpty(t, subdomain)
	origin := "https://" + subdomain + ".com"

	_, err = s.GetTrustedOriginID(ctx, origin)
	if err == nil {
		err = s.DeleteTrustedOrigin(ctx, origin)
		require.NoError(t, err)
		_, err = s.GetTrustedOriginID(ctx, origin)
		require.ErrorIs(t, err, ErrOriginNotFound)
	} else if !errors.Is(err, ErrOriginNotFound) {
		t.Fatalf("failed to get origin ID from the API: %s", err)
	}

	id, err := s.CreateTrustedOrigin(ctx, origin)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	t.Cleanup(func() {
		s.DeleteTrustedOrigin(ctx, id) //nolint:errcheck,gosec
	})

	err = s.DeleteTrustedOrigin(ctx, id)
	require.NoError(t, err)

	id, err = s.GetTrustedOriginID(ctx, origin)
	require.ErrorIs(t, err, ErrOriginNotFound)
	require.Empty(t, id)
}

func TestUpdateStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toRemove := []string{"1"}
		toAdd := []string{"4"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"2", "3", "4"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toRemove := []string{"3"}
		toAdd := []string{"4"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"1", "2", "4"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toRemove := []string{"1", "2", "3"}
		toAdd := []string{"3"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"3"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toAdd := []string{"3"}
		var toRemove []string
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"1", "2", "3"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		var toAdd []string
		var toRemove []string
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"1", "2", "3"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		var toAdd []string
		toRemove := []string{"4"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"1", "2", "3"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		var toAdd []string
		toRemove := []string{"2"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"1", "3"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		var source []string
		var toAdd []string
		toRemove := []string{"2"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		var source []string
		toRemove := []string{"2"}
		toAdd := []string{"1", "2"}
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{"1", "2"}, result)
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		var source []string
		var toRemove []string
		var toAdd []string
		result := UpdateStringsSet(source, toRemove, toAdd)
		require.Equal(t, []string{}, result)
	})
}

func TestValidateUpdateParams(t *testing.T) {
	t.Parallel()
	t.Run("invalid admin orgs", func(t *testing.T) {
		t.Parallel()
		params := UpdateUserParams{
			PortalAdminOrgs: &[]string{"aaa"},
			Lastname:        nil,
			Firstname:       nil,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, ErrInvalidPortalAdminOrgs, err)
	})

	t.Run("valid no first and last name", func(t *testing.T) {
		t.Parallel()
		params := UpdateUserParams{
			PortalAdminOrgs: &[]string{"ebe890cc-800f-11ec-a8a3-0242ac120002"},
			Lastname:        nil,
			Firstname:       nil,
		}
		err := validateUpdateUserParams(params)
		require.NoError(t, err)
	})

	t.Run("valid all", func(t *testing.T) {
		t.Parallel()
		firstName := "John"
		lastName := "Doe"
		params := UpdateUserParams{
			PortalAdminOrgs: &[]string{"ebe890cc-800f-11ec-a8a3-0242ac120002", "ebe890cc-800f-11ec-a8a3-0242ac120005"},
			PMMDemoIDs:      &[]string{"ebe890cc-800f-11ec-a8a3-0242ac120002", "ebe890cc-800f-11ec-a8a3-0242ac120005"},
			Lastname:        &lastName,
			Firstname:       &firstName,
		}
		err := validateUpdateUserParams(params)
		require.NoError(t, err)
	})

	t.Run("invalid PortalAdminOrgs not uuid", func(t *testing.T) {
		t.Parallel()
		firstName := "John"
		lastName := "Doe"
		params := UpdateUserParams{
			PortalAdminOrgs: &[]string{"aaa", "aaa"},
			Lastname:        &lastName,
			Firstname:       &firstName,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, err, ErrInvalidPortalAdminOrgs)
	})

	t.Run("invalid PortalAdminOrgs empty string", func(t *testing.T) {
		t.Parallel()
		params := UpdateUserParams{
			PortalAdminOrgs: &[]string{"ebe890cc-800f-11ec-a8a3-0242ac120002", ""},
			Lastname:        nil,
			Firstname:       nil,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, err, ErrInvalidPortalAdminOrgs)
	})

	t.Run("empty last name", func(t *testing.T) {
		t.Parallel()
		lastName := ""
		params := UpdateUserParams{
			PortalAdminOrgs: nil,
			Lastname:        &lastName,
			Firstname:       nil,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, err, ErrEmptyLastName)
	})

	t.Run("empty first name", func(t *testing.T) {
		t.Parallel()
		firstName := ""
		params := UpdateUserParams{
			PortalAdminOrgs: nil,
			Lastname:        nil,
			Firstname:       &firstName,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, err, ErrEmptyFirstName)
	})

	t.Run("invalid PMM demo id", func(t *testing.T) {
		t.Parallel()
		params := UpdateUserParams{
			PMMDemoIDs: &[]string{"ebe890cc-800f-11ec-a8a3-0242ac120002", ""},
			Lastname:   nil,
			Firstname:  nil,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, err, ErrInvalidPMMDemoID)
	})

	t.Run("duplicated pmm demo ID", func(t *testing.T) {
		t.Parallel()
		params := UpdateUserParams{
			PMMDemoIDs: &[]string{"ebe890cc-800f-11ec-a8a3-0242ac120002", "ebe890cc-800f-11ec-a8a3-0242ac120002"},
			Lastname:   nil,
			Firstname:  nil,
		}
		err := validateUpdateUserParams(params)
		require.ErrorIs(t, err, ErrDuplicatedPMMDemoID)
	})
}

func TestCreateOAuthApp(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	ctx := context.Background()
	pmmServerID := gofakeit.UUID()

	params := &OAuthAppParams{
		PMMServerID:          pmmServerID,
		PMMServerURL:         "https://localhost/graph",
		PMMServerCallbackURL: "https://localhost/graph/login/generic_oauth",
		OrgID:                uuid.NewString(),
		InventoryID:          pmmServerID,
	}

	app, err := s.CreateOAuthApp(ctx, params)
	require.NoError(t, err)
	t.Cleanup(func() {
		s.DeleteApp(ctx, app.AppID) //nolint:errcheck,gosec
	})
	require.NotNil(t, app)
	require.NotEmpty(t, app.Credentials.OAuthClient.ClientID)
}

func TestCreateMachineAuthApp(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	ctx := context.Background()
	pmmServerID := uuid.NewString()

	params := &MachineAuthAppParams{
		PMMServerID: pmmServerID,
		OrgID:       uuid.NewString(),
		InventoryID: pmmServerID,
	}

	app, err := s.CreateMachineAuthApp(ctx, params)
	require.NoError(t, err)
	t.Cleanup(func() {
		s.DeleteApp(ctx, app.AppID) //nolint:errcheck,gosec
	})
	require.NotNil(t, app)
	require.NotEmpty(t, app.Credentials.OAuthClient.ClientID)
	require.NotEmpty(t, app.Credentials.OAuthClient.ClientSecret)
}

func TestGetActivationLink(t *testing.T) {
	t.Parallel()

	t.Run("not activated user", func(t *testing.T) {
		t.Parallel()

		s, err := createOktaService(t)
		require.NoError(t, err)

		ctx := context.Background()

		email, password, firstName, lastName := GenCredentials(t)
		testUser := CreateInactivatedTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, testUser.ID)
		})

		info, err := s.GetActivationInfo(ctx, testUser.ID)
		require.NoError(t, err)
		require.NotEmpty(t, info)
	})

	t.Run("activated user", func(t *testing.T) {
		t.Parallel()

		s, err := createOktaService(t)
		require.NoError(t, err)

		ctx := context.Background()

		email, password, firstName, lastName := GenCredentials(t)
		testUser := CreateTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, testUser.ID)
		})

		info, err := s.GetActivationInfo(ctx, testUser.ID)
		require.Error(t, err)
		require.Empty(t, info)
	})
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	t.Run("valid values", func(t *testing.T) {
		t.Parallel()

		login := "test"

		profile := okta.UserProfile{
			profileLogin:     login,
			profileMarketing: true,
		}

		valueStr, err := getValue[string](profile, profileLogin)
		require.NoError(t, err)
		require.Equal(t, login, *valueStr)

		valueBool, err := getValue[bool](profile, profileMarketing)
		require.NoError(t, err)
		require.True(t, *valueBool)

		valueStr, err = getValue[string](profile, profileFirstName)
		require.ErrorIs(t, err, errNotFound)
		require.Nil(t, valueStr)
	})

	t.Run("wrong field format", func(t *testing.T) {
		t.Parallel()

		profile := okta.UserProfile{
			profileLogin: []string{},
		}

		valueStr, err := getValue[string](profile, profileLogin)
		require.EqualError(t, err, "unexpected field type")
		require.Nil(t, valueStr)

		valueBool, err := getValue[bool](profile, profileLogin)
		require.EqualError(t, err, "unexpected field type")
		require.Nil(t, valueBool)
	})
}

func TestGetReactivationInfo(t *testing.T) {
	t.Parallel()

	t.Run("error: user status ACTIVE", func(t *testing.T) {
		t.Parallel()

		s, err := createOktaService(t)
		require.NoError(t, err)

		ctx := context.Background()

		email, password, firstName, lastName := GenCredentials(t)
		testUser := CreateTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, testUser.ID)
		})
		require.Equal(t, UserStatusActive, testUser.Status)

		link, err := s.GetReactivationInfo(ctx, testUser.ID)
		require.ErrorContains(t, err, "This operation is not allowed in the user's current status.")
		require.Empty(t, link)
	})

	t.Run("success: user status PROVISIONED", func(t *testing.T) {
		t.Parallel()

		s, err := createOktaService(t)
		require.NoError(t, err)

		ctx := context.Background()

		email, _, _, _ := GenCredentials(t)

		u := okta.CreateUserRequest{ //nolint:exhaustivestruct
			Profile: &okta.UserProfile{
				profileEmail: email,
				profileLogin: email,
			},
		}
		qp := query.NewQueryParams(query.WithActivate(false))
		user, _, err := createOktaClient(t).User.CreateUser(context.Background(), u, qp)
		require.NoError(t, err)

		t.Cleanup(func() {
			DeleteUser(t, user.Id)
		})
		require.Equal(t, UserStatusStaged, user.Status)

		activationInfo, err := s.GetActivationInfo(ctx, user.Id)
		require.NoError(t, err)
		require.NotEmpty(t, activationInfo)

		updatedUser, _, err := s.c.User.GetUser(ctx, user.Id)
		require.Equal(t, UserStatusProvisioned, updatedUser.Status)
		require.NoError(t, err)

		link, err := s.GetReactivationInfo(ctx, user.Id)
		require.NoError(t, err)
		require.NotEmpty(t, link)

		user, _, err = s.c.User.GetUser(ctx, user.Id)
		require.Equal(t, UserStatusProvisioned, user.Status)
		require.NoError(t, err)
	})

	t.Run("error: user status STAGED", func(t *testing.T) {
		t.Parallel()

		s, err := createOktaService(t)
		require.NoError(t, err)

		ctx := context.Background()

		email, password, firstName, lastName := GenCredentials(t)
		testUser := CreateInactivatedTestUser(t, email, password, firstName, lastName)
		t.Cleanup(func() {
			DeleteUser(t, testUser.ID)
		})

		require.Equal(t, UserStatusStaged, testUser.Status)

		link, err := s.GetReactivationInfo(ctx, testUser.ID)
		require.ErrorContains(t, err, "This operation is not allowed in the user's current status.")
		require.Empty(t, link)
	})
}

func TestGetAppSecret(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	ctx := context.Background()

	params := new(MachineAuthAppParams)
	// create some app
	app, err := s.CreateMachineAuthApp(ctx, params)
	require.NoError(t, err)
	t.Cleanup(func() {
		s.DeleteApp(ctx, app.AppID) //nolint:errcheck,gosec
	})

	secret, err := s.GetAppSecret(ctx, app.AppID)
	require.NotEmpty(t, secret)
	// check the secret is not empty
	require.NotEmpty(t, *secret)
	require.NoError(t, err)
}
