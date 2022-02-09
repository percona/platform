package okta

const (
	// StatusStaged Accounts have a staged status when they are first created, before the activation flow is initiated,
	// or if there is a pending admin action.
	StatusStaged = "STAGED"

	// StatusProvisioned 	Accounts have a provisioned status when they are provisioned, but the user has not provided
	// verification by clicking through the activation email or provided a password.
	StatusProvisioned = "PROVISIONED"

	// StatusActive 	Accounts have an active status when:
	//   - An admin adds a user (Add person) on the People page and sets the user password without requiring email verification.
	//   - An admin adds a user (Add person) on the People page, sets the user password, and requires the user to set their
	//     password when they first sign-in.
	//   - A user self-registers into a custom app or the Okta Homepage and email verification is not required.
	//   - An admin explicitly activate user accounts.
	StatusActive = "ACTIVE"

	// StatusRecovery Accounts have a recovery status when a user requests a password reset or an admin initiates one on their behalf.
	StatusRecovery = "RECOVERY"

	// StatusPasswordExpired Accounts have a password expired status when the password has expired and the account requires an
	// update to the password before a user is granted access to applications.
	StatusPasswordExpired = "PASSWORD EXPIRED"

	// StatusLockedOut Accounts have a locked out status when the user exceeds the number of login attempts defined in the login policy.
	StatusLockedOut = "LOCKED OUT"

	// StatusSuspended Accounts have a suspended status when an admin explicitly suspends them. The user cannot access applications,
	// the Admin Console, or the Okta End-User Dashboard. Application assignments are unaffected and the user profile can be updated.
	StatusSuspended = "SUSPENDED"

	// StatusDeprovisioned Accounts have a deprovisioned status when an admin explicitly deactivates or deprovisions them.
	// All application assignments are removed and the password is permanently deleted.
	StatusDeprovisioned = "DEPROVISIONED"
)
