package models

// SetID represents a globally unique identifier for this Set
type SetID int

// Set tracks a single match in a bracket
type Set struct {
	ID       SetID
	WinnerID ParticipantID
	Slots    []Slot
}
