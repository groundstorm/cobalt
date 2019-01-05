package models

// MatchSlotID represents a globally unique identifier for this Match
type MatchSlotID string

// MatchSlot represents either P1 or P2 in a match
type MatchSlot struct {
	ID     MatchSlotID
	UserID UserID
}
