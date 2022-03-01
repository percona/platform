package okta

// contains custom types returned by Okta client code.

// User represents user structure.
type User struct {
	ID              string
	Login           string
	FirstName       string
	LastName        string
	Status          string
	PortalAdminOrgs []string
}

// UpdateUserParams parameters set to update a user.
type UpdateUserParams struct {
	PortalAdminOrgs *[]string
	Lastname        *string
	Firstname       *string
}

// RegisterUserParams parameters set to invite a user.
type RegisterUserParams struct {
	Login string
}

// Group represents user group structure.
type Group struct {
	ID          string
	Name        string
	Description string
}
