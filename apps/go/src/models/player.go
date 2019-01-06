package models

// PlayerID represents a globally unique identifier for this user
type PlayerID int

// The Email represents a user's email address
type Email string

// The Player struct represents any user of the system.  Players, spectators,
// TOs, Judges, etc.  All must have a user account to do anything
type Player struct {
	ID        PlayerID `json:"id"`
	Email     Email    `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
}
