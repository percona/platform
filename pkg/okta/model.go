package okta

// contains custom types returned by Okta client code.

// User represents user structure.
type User struct {
	ID        string
	Login     string
	FirstName string
	LastName  string
	Status    string
}

// Group represents user group structure.
type Group struct {
	ID          string
	Name        string
	Description string
}
