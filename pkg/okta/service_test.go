package okta

import (
	"context"
	"os"
	"os/user"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const devHost = "okta-dev.percona.com"

var authErrorType = &AuthError{} //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	gofakeit.Seed(time.Now().UnixNano())
}

func getOktaToken(t *testing.T) string {
	token := os.Getenv("OKTA_TOKEN")
	require.NotEmpty(t, token, "Okta token is missing")

	return token
}

// createTestUser signs up an okta user with a password unlike our registration flow
// since we need a user with a set password in our tests.
func createTestUser(t *testing.T, s *Service, email, password, firstName, lastName string) *User {
	t.Helper()

	u := okta.CreateUserRequest{ //nolint:exhaustivestruct
		Profile: &okta.UserProfile{
			"email":     email,
			"login":     email,
			"firstName": firstName,
			"lastName":  lastName,
		},
		Credentials: &okta.UserCredentials{
			Password: &okta.PasswordCredential{
				Value: password,
			},
		},
	}
	qp := query.NewQueryParams(query.WithActivate(true))
	user, _, err := s.c.User.CreateUser(context.Background(), u, qp)
	require.NoError(t, err)

	nLogin, err := getUserLogin(user)
	require.NoError(t, err)

	return &User{
		ID:     user.Id,
		Login:  nLogin,
		Status: user.Status,
	}
}

func TestSignUp(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	t.Run("invalid login", func(t *testing.T) {
		_, _, firstName, lastName := genCredentials(t)
		user, err := s.SignUp(context.Background(), "not email", firstName, lastName)
		require.EqualError(t, err, "invalid login: login: Username must be in the form of an email address")
		require.IsType(t, authErrorType, err)
		assert.Nil(t, user)
	})

	t.Run("empty login", func(t *testing.T) {
		_, _, firstName, lastName := genCredentials(t)
		user, err := s.SignUp(context.Background(), "", firstName, lastName)
		require.Equal(t, err, ErrEmptyLogin)
		require.IsType(t, authErrorType, err)
		assert.Nil(t, user)
	})

	t.Run("empty first name", func(t *testing.T) {
		email, _, _, lastName := genCredentials(t)
		user, err := s.SignUp(context.Background(), email, "", lastName)
		require.Equal(t, err, ErrEmptyFirstName)
		require.IsType(t, authErrorType, err)
		assert.Nil(t, user)
	})

	t.Run("empty last name", func(t *testing.T) {
		email, _, firstName, _ := genCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, "")
		require.Equal(t, err, ErrEmptyLastName)
		require.IsType(t, authErrorType, err)
		assert.Nil(t, user)
	})

	t.Run("valid sign up", func(t *testing.T) {
		email, _, firstName, lastName := genCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, lastName)
		require.NoError(t, err)
		defer deleteUser(t, s, user.ID)

		assert.Equal(t, email, user.Login)
		assert.Equal(t, user.Status, "PROVISIONED")
		assert.NotEmpty(t, user.ID)
	})
}

func TestSignIn(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	user := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, user.ID)

	t.Run("invalid password", func(t *testing.T) {
		userID, sessionToken, err := s.SignIn(context.Background(), email, "wrong")
		require.Equal(t, ErrAuthentication, err)
		require.IsType(t, authErrorType, err)
		assert.Empty(t, sessionToken)
		assert.Empty(t, userID)
	})

	t.Run("empty password", func(t *testing.T) {
		userID, sessionToken, err := s.SignIn(context.Background(), email, "")
		require.Equal(t, ErrEmptyPassword, err)
		require.IsType(t, authErrorType, err)
		assert.Empty(t, sessionToken)
		assert.Empty(t, userID)
	})

	t.Run("invalid login", func(t *testing.T) {
		userID, sessionToken, err := s.SignIn(context.Background(), "wrong", password)
		require.Equal(t, ErrAuthentication, err)
		require.IsType(t, authErrorType, err)
		assert.Empty(t, sessionToken)
		assert.Empty(t, userID)
	})

	t.Run("empty login", func(t *testing.T) {
		userID, sessionToken, err := s.SignIn(context.Background(), "", password)
		require.Equal(t, err, ErrEmptyLogin)
		require.IsType(t, authErrorType, err)
		assert.Empty(t, sessionToken)
		assert.Empty(t, userID)
	})

	t.Run("valid sign in", func(t *testing.T) {
		userID, sessionToken, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)
		assert.NotEmpty(t, sessionToken)
		assert.NotEmpty(t, userID)
	})
}

func TestSessions(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	t.Run("invalid session", func(t *testing.T) {
		login, err := s.CheckSession(context.Background(), "invalid-session-token")
		require.Equal(t, err, ErrNotFound)
		require.IsType(t, authErrorType, err)
		assert.Empty(t, login)
	})

	t.Run("valid session", func(t *testing.T) {
		email, password, firstName, lastName := genCredentials(t)
		user := createTestUser(t, s, email, password, firstName, lastName)
		defer deleteUser(t, s, user.ID)

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
		assert.GreaterOrEqual(t, expiresAt.Unix(), ts.Add(sessionTTL-timeError).Unix())
		assert.LessOrEqual(t, expiresAt.Unix(), time.Now().Add(sessionTTL+timeError).Unix())

		userEmail, err := s.CheckSession(context.Background(), sessionID)
		require.NoError(t, err)
		assert.Equal(t, email, userEmail)
	})
}

func TestSessionRefresh(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	user := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, user.ID)

	t.Run("normal", func(t *testing.T) {
		_, token, err := s.SignIn(context.Background(), email, password)
		require.NoError(t, err)

		sessionID, expiresAt, err := s.CreateSession(context.Background(), token)
		require.NoError(t, err)

		time.Sleep(time.Second)

		newExpirationTime, err := s.RefreshSession(context.Background(), sessionID)
		require.NoError(t, err)

		assert.Greater(t, newExpirationTime.Unix(), expiresAt.Unix())
	})

	t.Run("invalid session", func(t *testing.T) {
		expTime, err := s.RefreshSession(context.Background(), "invalid-session-id")
		require.Equal(t, err, ErrNotFound)
		require.Zero(t, expTime)
	})
}

func TestCloseSession(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	user := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, user.ID)

	t.Run("normal", func(t *testing.T) {
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
		err = s.CloseSession(context.Background(), "invalid-session-id")
		require.Equal(t, err, ErrNotFound)
	})

	t.Run("already closed session", func(t *testing.T) {
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
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	user := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, user.ID)

	t.Run("user doesn't exists", func(t *testing.T) {
		userID, err := s.FindUser(context.Background(), "invalid@example.com")
		require.Equal(t, ErrNotFound, err)
		require.Empty(t, userID)
	})

	t.Run("user exists", func(t *testing.T) {
		userID, err := s.FindUser(context.Background(), email)
		require.NoError(t, err)
		require.NotEmpty(t, userID)
	})
}

func TestPasswordReset(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	user := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, user.ID)

	u, err := s.FindUser(context.Background(), email)
	require.NoError(t, err)

	err = s.ResetPassword(context.Background(), u.Id)
	assert.NoError(t, err)

	_, _, err = s.SignIn(context.Background(), email, password)
	assert.Equal(t, ErrAuthentication, err)
}

func TestGroups(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	user := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, user.ID)

	t.Run("create group and add user", func(t *testing.T) {
		name := gofakeit.UUID()
		description := "Test group"
		group, err := s.CreateGroup(context.Background(), name, description)
		defer deleteGroup(t, s, group.ID)
		require.NoError(t, err)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, description, group.Description)
		assert.NotEmpty(t, group.ID)

		err = s.AddUserToGroup(context.Background(), user.ID, group.ID)
		require.NoError(t, err)

		users, err := s.GetGroupMembers(context.Background(), group.ID, 0, "")
		require.NoError(t, err)

		assert.Len(t, users, 1)
		assert.Equal(t, *user, users[0])
	})
}

func TestGetUserLogin(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		login := "test"
		user := okta.User{
			Profile: &okta.UserProfile{
				"login": login,
			},
		}

		actual, err := getUserLogin(&user)
		require.NoError(t, err)
		assert.Equal(t, login, actual)
	})

	t.Run("missing login", func(t *testing.T) {
		user := okta.User{
			Profile: &okta.UserProfile{},
		}

		login, err := getUserLogin(&user)
		require.EqualError(t, err, "missing user login")
		assert.Empty(t, login)
	})

	t.Run("missing user profile", func(t *testing.T) {
		var user okta.User

		login, err := getUserLogin(&user)
		require.EqualError(t, err, "missing user profile")
		assert.Empty(t, login)
	})
}

func TestDeleteUser(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	t.Run("valid", func(t *testing.T) {
		email, _, firstName, lastName := genCredentials(t)
		user, err := s.SignUp(context.Background(), email, firstName, lastName)
		require.NoError(t, err)

		err = s.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)

		_, err = s.FindUser(context.Background(), user.Login)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("missing user", func(t *testing.T) {
		err = s.DeleteUser(context.Background(), "unknown-id")
		require.Equal(t, ErrNotFound, err)
	})
}

func TestUpdateProfile(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	email, password, firstName, lastName := genCredentials(t)
	testUser := createTestUser(t, s, email, password, firstName, lastName)
	defer deleteUser(t, s, testUser.ID)

	t.Run("user doesn't exists", func(t *testing.T) {
		_, err := s.UpdateProfile(context.Background(), &okta.User{Id: "unknown", Profile: &okta.UserProfile{}}, "firstName", "lastName")
		require.EqualError(t, err, "missing user login")
	})

	t.Run("user exists update successful", func(t *testing.T) {
		user, err := s.FindUser(context.Background(), testUser.Login)
		require.NoError(t, err)
		prof := *user.Profile
		t.Log(prof["firstName"], prof["lastName"], prof["login"], user.Id)

		newFirstName := gofakeit.FirstName()
		newLastName := gofakeit.LastName()
		updatedUser, err := s.UpdateProfile(context.Background(), user, newFirstName, newLastName)
		require.NoError(t, err)

		updatedProfile := *updatedUser.Profile
		assert.Equal(t, user.Id, updatedUser.Id)
		assert.Equal(t, newFirstName, updatedProfile["firstName"])
		assert.Equal(t, newLastName, updatedProfile["lastName"])
	})
}

func TestGetRegisteredUsersCount(t *testing.T) {
	s, err := New(devHost, getOktaToken(t))
	require.NoError(t, err)

	usersCount, err := s.GetRegisteredUsersCount(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, usersCount)
}

// genCredentials creates test user email and password.
func genCredentials(t *testing.T) (string, string, string, string) {
	t.Helper()

	hostname, err := os.Hostname()
	require.NoError(t, err)

	u, err := user.Current()
	require.NoError(t, err)

	email := strings.Join([]string{u.Username, hostname, gofakeit.Email(), "test"}, ".")
	password := gofakeit.Password(true, true, true, false, false, 14)
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	return email, password, firstName, lastName
}

func deleteUser(t *testing.T, s *Service, userID string) {
	t.Helper()

	err := s.DeleteUser(context.Background(), userID)
	assert.NoError(t, err)
}

func deleteGroup(t *testing.T, s *Service, groupID string) {
	t.Helper()

	_, err := s.c.Group.DeleteGroup(context.Background(), groupID)
	assert.NoError(t, err)
}
