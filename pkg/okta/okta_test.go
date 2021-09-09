package okta

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/require"

	"github.com/percona-platform/platform/pkg/model"
	"github.com/percona-platform/platform/pkg/testutils"
)

var authErrorType = &AuthError{} //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	gofakeit.Seed(time.Now().UnixNano())
}

func createOktaService(t *testing.T) (*Service, error) {
	t.Helper()
	return New(context.Background(), testutils.OktaDevHost, testutils.GetOktaToken(t))
}

func TestSignUp(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	t.Run("invalid login", func(t *testing.T) {
		t.Parallel()

		_, _, firstName, lastName := testutils.GenCredentials(t)
		user, err := s.SignUp(context.Background(), "not email", firstName, lastName)
		require.EqualError(t, err, "invalid login: login: Username must be in the form of an email address")
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("empty login", func(t *testing.T) {
		t.Parallel()

		_, _, firstName, lastName := testutils.GenCredentials(t)
		user, err := s.SignUp(context.Background(), "", firstName, lastName)
		require.Equal(t, err, ErrEmptyLogin)
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("empty first name", func(t *testing.T) {
		t.Parallel()

		email, _, _, lastName := testutils.GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, "", lastName)
		require.Equal(t, err, ErrEmptyFirstName)
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("empty last name", func(t *testing.T) {
		t.Parallel()

		email, _, firstName, _ := testutils.GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, "")
		require.Equal(t, err, ErrEmptyLastName)
		require.IsType(t, authErrorType, err)
		require.Nil(t, user)
	})

	t.Run("valid sign up", func(t *testing.T) {
		t.Parallel()

		email, _, firstName, lastName := testutils.GenCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, lastName)
		require.NoError(t, err)
		defer testutils.DeleteUser(t, user.ID)

		require.Equal(t, email, user.Login)
		require.Equal(t, user.Status, "PROVISIONED")
		require.NotEmpty(t, user.ID)
	})
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
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

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
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

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
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

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
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

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
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

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
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

	email, password, firstName, lastName := testutils.GenCredentials(t)
	user := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, user.ID)
	})

	name := gofakeit.UUID()
	description := "Test group"
	group, err := s.CreateGroup(context.Background(), name, description)
	t.Cleanup(func() {
		testutils.DeleteGroup(t, group.ID)
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
			Profile: &okta.UserProfile{},
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

		email, _, firstName, lastName := testutils.GenCredentials(t)
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

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	s, err := createOktaService(t)
	require.NoError(t, err)

	email, password, firstName, lastName := testutils.GenCredentials(t)
	testUser := testutils.CreateTestUser(t, email, password, firstName, lastName)
	t.Cleanup(func() {
		testutils.DeleteUser(t, testUser.ID)
	})

	t.Run("user doesn't exists", func(t *testing.T) {
		t.Parallel()

		_, err := s.UpdateProfile(context.Background(), &model.User{ID: "unknown", Login: "login", Status: "status"}, "firstName", "lastName")
		require.EqualError(t, err, "missing user login")
	})

	t.Run("user exists update successful", func(t *testing.T) {
		t.Parallel()

		user, err := s.FindUser(context.Background(), testUser.Login)
		require.NoError(t, err)
		t.Log(user.FirstName, user.LastName, user.Login, user.ID)

		newFirstName := gofakeit.FirstName()
		newLastName := gofakeit.LastName()
		updatedUser, err := s.UpdateProfile(context.Background(), user, newFirstName, newLastName)
		require.NoError(t, err)

		require.Equal(t, user.ID, updatedUser.ID)
		require.Equal(t, newFirstName, updatedUser.FirstName)
		require.Equal(t, newLastName, updatedUser.LastName)
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
