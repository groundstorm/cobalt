package models

// ID represents a globally unique identifier for this Group
type GroupID int

// Group tracks a single bracket in a tournament
type Group struct {
	ID      GroupID
	PhaseID PhaseID
	Sets    []*Set
}
