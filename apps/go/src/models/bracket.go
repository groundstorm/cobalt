package models

// ID represents a globally unique identifier for this Bracket
type BracketID string

// Bracket tracks a single bracket in a tournament
type Bracket struct {
	ID BracketID
}
