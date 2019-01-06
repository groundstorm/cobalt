package models

// SlotID represents a globally unique identifier for this Set
type SlotID int

// Slot represents either P1 or P2 in a match
type Slot struct {
	ID            SlotID
	ParticipantID ParticipantID
	SeedID        SeedID
}
