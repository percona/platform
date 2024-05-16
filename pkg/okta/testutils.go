package okta

import (
	"context"
	"net/url"
	"os"
	"os/user"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
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
	return createTestUser(t, email, password, firstName, lastName, true)
}

// CreateInactivatedTestUser signs up an okta user with a password but not activates them.
func CreateInactivatedTestUser(t *testing.T, email, password, firstName, lastName string) *User {
	t.Helper()
	return createTestUser(t, email, password, firstName, lastName, false)
}

// ActivateUser activates an existing user.
func ActivateUser(t *testing.T, id string) string {
	t.Helper()

	qp := query.NewQueryParams(query.WithSendEmail(false))
	tokenInfo, _, err := createOktaClient(t).User.ActivateUser(context.Background(), id, qp)
	require.NoError(t, err)
	require.NotEmpty(t, tokenInfo)

	return tokenInfo.ActivationToken
}

func createTestUser(t *testing.T, email, password, firstName, lastName string, activate bool) *User {
	t.Helper()

	u := okta.CreateUserRequest{ //nolint:exhaustivestruct
		Profile: &okta.UserProfile{
			profileEmail:           email,
			profileLogin:           email,
			profileFirstName:       firstName,
			profileLastName:        lastName,
			profilePortalAdminOrgs: []string{},
			profileSecondaryEmail:  gofakeit.Email(),
			profileMobilePhone:     gofakeit.Phone(),
		},
	}
	if password != "" {
		u.Credentials = &okta.UserCredentials{
			Password: &okta.PasswordCredential{
				Value: password,
			},
		}
	}

	qp := query.NewQueryParams(query.WithActivate(activate))
	testUser, _, err := createOktaClient(t).User.CreateUser(context.Background(), u, qp)
	require.NoError(t, err, "failed to create test user with password: "+password)

	converterUser, err := convertUser(testUser)
	require.NoError(t, err)

	return converterUser
}

// GenCredentials create test user email and password.
func GenCredentials(t *testing.T) (string, string, string, string) {
	t.Helper()

	hostname, err := os.Hostname()
	require.NoError(t, err)

	u, err := user.Current()
	require.NoError(t, err)

	email := strings.Join([]string{u.Username, hostname, gofakeit.Email(), "test"}, ".")
	password := GenPassword(t)
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	return email, password, firstName, lastName
}

// GenPassword generates a password with at least one lowercase, uppercase, special character, and digit.
func GenPassword(t *testing.T) string {
	t.Helper()

	const lowerStr = "abcdefghijklmnopqrstuvwxyz"
	const upperStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const specialSafeStr = "!@#$%^&*"
	const numericStr = "0123456789"

	faker := gofakeit.New(time.Now().UnixNano())
	// All params set to true in faker.Password() call are optional
	// and do not guarantee that password will contain requested characters.
	password := faker.Password(true, true, true, true, false, 14) //nolint:mnd
	// make sure that password contains at least one lowercase character.
	password += string(lowerStr[faker.Rand.Int63()%int64(len(lowerStr))])
	// make sure that password contains at least one uppercase character.
	password += string(upperStr[faker.Rand.Int63()%int64(len(upperStr))])
	// make sure that password contains at least one special character.
	password += string(specialSafeStr[faker.Rand.Int63()%int64(len(specialSafeStr))])
	// make sure that password contains at least one digit.
	password += string(numericStr[faker.Rand.Int63()%int64(len(numericStr))])
	return password
}

// DeleteUser delete user from Okta by UserID.
func DeleteUser(t *testing.T, userID string) {
	t.Helper()

	_, err := createOktaClient(t).User.DeactivateOrDeleteUser(context.Background(), userID, nil)
	require.NoError(t, err)
}

// DeleteGroup delete group from Okta by GroupID.
func DeleteGroup(t *testing.T, groupID string) {
	t.Helper()

	_, err := createOktaClient(t).Group.DeleteGroup(context.Background(), groupID)
	require.NoError(t, err)
}

func oktaAPIRequest(oktaClient *okta.Client, method, path string, body, v interface{}) error {
	requestExecutor := oktaClient.CloneRequestExecutor().WithAccept("application/json").WithContentType("application/json")
	req, err := requestExecutor.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	resp, err := requestExecutor.Do(context.Background(), req, v)
	if err != nil {
		return err
	}
	defer closeResponseBody(resp)

	return err
}
