package models

// ID represents a globally unique identifier for this Match
type MatchID string

// Match tracks a single match in a bracket
type Match struct {
	ID MatchID
}
