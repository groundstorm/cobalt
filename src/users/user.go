package users

// UserID represents a globally unique identifier for this user
type UserID string

// The UserEmail represents a user's email address
type UserEmail string

// The User struct represents any user of the system.  Players, spectators,
// TOs, Judges, etc.  All must have a user account to do anything
type User struct {
	ID        UserID
	Email     UserEmail
	FirstName string
	LastName  string
}
