package models

// EntrantID represents a globally unique identifier for this user
type EntrantID int

// The Entrant struct represents any user of the system.  Players, spectators,
// TOs, Judges, etc.  All must have a user account to do anything
type Entrant struct {
	ID            EntrantID
	ParticipantID ParticipantID
}
