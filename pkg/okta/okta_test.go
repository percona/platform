package okta

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
		require.Equal(t, user.Status, "PROVISIONED")
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

		time.Sleep(time.Second)

		newExpirationTime, err := s.RefreshSession(context.Background(), sessionID)
		require.NoError(t, err)

		require.Greater(t, newExpirationTime.Unix(), expiresAt.Unix())
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

		userID, err := s.FindUser(context.Background(), email)
		require.NoError(t, err)
		require.NotEmpty(t, userID)
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
	assert.True(t, exists)

	exists, err = s.GroupExists(context.Background(), "non-existent-group")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestGetUserLogin(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		login := "test"
		user := okta.User{
			Profile: &okta.UserProfile{
				"login": login,
			},
		}

		actual, err := getUserLogin(&user)
		require.NoError(t, err)
		require.Equal(t, login, actual)
	})

	t.Run("missing login", func(t *testing.T) {
		t.Parallel()

		user := okta.User{
			Profile: new(okta.UserProfile),
		}

		login, err := getUserLogin(&user)
		require.EqualError(t, err, "missing user login")
		require.Empty(t, login)
	})

	t.Run("missing user profile", func(t *testing.T) {
		t.Parallel()

		var user okta.User

		login, err := getUserLogin(&user)
		require.EqualError(t, err, "missing user profile")
		require.Empty(t, login)
	})
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

		_, err := s.UpdateUser(context.Background(), "unknown", UpdateUserParams{Firstname: "firstName", Lastname: "lastName"})
		require.EqualError(t, err, "not found")
	})

	t.Run("user exists update lastname firstname successful", func(t *testing.T) {
		t.Parallel()

		user, err := s.FindUser(context.Background(), testUser.Login)
		require.NoError(t, err)
		t.Log(user.FirstName, user.LastName, user.Login, user.ID)

		newFirstName := gofakeit.FirstName()
		newLastName := gofakeit.LastName()
		updatedUser, err := s.UpdateUser(context.Background(), user.ID, UpdateUserParams{Firstname: newFirstName, Lastname: newLastName})
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
		updatedUser, err := s.UpdateUser(context.Background(), user.ID, UpdateUserParams{PortalAdminOrgsToAdd: []string{newID}})
		require.NoError(t, err)

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
	assert.NotEmpty(t, app.AppID)
	require.NotNil(t, app.Credentials)
	require.NotNil(t, app.Credentials.OAuthClient)
	assert.NotEmpty(t, app.Credentials.OAuthClient.ClientID)
	assert.NotEmpty(t, app.Credentials.OAuthClient.ClientSecret)

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
	assert.False(t, assigned)

	err = s.AddAppToGroup(ctx, app.AppID, group.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		s.RemoveAppFromGroup(ctx, app.AppID, group.ID) //nolint:errcheck,gosec
	})

	assigned = s.IsAppAssignedToGroup(ctx, app.AppID, group.ID)
	assert.True(t, assigned)

	err = s.RemoveAppFromGroup(ctx, app.AppID, group.ID)
	require.NoError(t, err)

	assigned = s.IsAppAssignedToGroup(ctx, app.AppID, group.ID)
	assert.False(t, assigned)

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
		assert.ErrorIs(t, err, ErrOriginNotFound)
	} else if !errors.Is(err, ErrOriginNotFound) {
		t.Fatalf("failed to get origin ID from the API: %s", err)
	}

	id, err := s.CreateTrustedOrigin(ctx, origin)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	t.Cleanup(func() {
		s.DeleteTrustedOrigin(ctx, id) //nolint:errcheck,gosec
	})

	err = s.DeleteTrustedOrigin(ctx, id)
	require.NoError(t, err)

	id, err = s.GetTrustedOriginID(ctx, origin)
	require.ErrorIs(t, err, ErrOriginNotFound)
	assert.Empty(t, id)
}

func TestUpdateStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toRemove := []string{"1"}
		toAdd := []string{"4"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"2", "3", "4"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toRemove := []string{"3"}
		toAdd := []string{"4"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"1", "2", "4"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toRemove := []string{"1", "2", "3"}
		toAdd := []string{"3"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"3"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		toAdd := []string{"3"}
		var toRemove []string
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"1", "2", "3"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		var toAdd []string
		var toRemove []string
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"1", "2", "3"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		var toAdd []string
		toRemove := []string{"4"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"1", "2", "3"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		source := []string{"1", "2", "3"}
		var toAdd []string
		toRemove := []string{"2"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"1", "3"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		var source []string
		var toAdd []string
		toRemove := []string{"2"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		var source []string
		toRemove := []string{"2"}
		toAdd := []string{"1", "2"}
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{"1", "2"})
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		var source []string
		var toRemove []string
		var toAdd []string
		result := updateStringSlice(source, toRemove, toAdd)
		assert.Equal(t, result, []string{})
	})
}
