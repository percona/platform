package okta

import (
	"context"
	"net/url"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// OktaDevHost Okta domain.
const OktaDevHost = "id-dev.percona.com"

// GetOktaToken reads OKTA_TOKEN env var and return its value.
func GetOktaToken(t *testing.T) string {
	t.Helper()

	token := os.Getenv("OKTA_TOKEN")
	require.NotEmpty(t, token, "OKTA_TOKEN env var is missing")

	return token
}

// createOktaClient creates Okta client.
func createOktaClient(t *testing.T) *okta.Client {
	t.Helper()

	u := url.URL{Scheme: "https", Host: OktaDevHost}

	_, client, _ := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(u.String()),
		okta.WithToken(GetOktaToken(t)),
		okta.WithCache(false),
	)

	return client
}

// CreateTestUser signs up an okta user with a password unlike our registration flow
// since we need a user with a set password in our tests.
func CreateTestUser(t *testing.T, email, password, firstName, lastName string) *User {
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
	testUser, _, err := createOktaClient(t).User.CreateUser(context.Background(), u, qp)
	require.NoError(t, err)

	nLogin, err := getUserLogin(testUser)
	require.NoError(t, err)

	return &User{
		ID:     testUser.Id,
		Login:  nLogin,
		Status: testUser.Status,
	}
}

// GenCredentials create test user email and password.
func GenCredentials(t *testing.T) (string, string, string, string) {
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

// DeleteUser delete user from Okta by UserID.
func DeleteUser(t *testing.T, userID string) {
	t.Helper()

	_, err := createOktaClient(t).User.DeactivateOrDeleteUser(context.Background(), userID, nil)
	assert.NoError(t, err)
}

// DeleteGroup delete group from Okta by GroupID.
func DeleteGroup(t *testing.T, groupID string) {
	t.Helper()

	_, err := createOktaClient(t).Group.DeleteGroup(context.Background(), groupID)
	assert.NoError(t, err)
}
