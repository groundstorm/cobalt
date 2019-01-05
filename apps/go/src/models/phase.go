package models

// PhaseID represents a globally unique identifier for this event
type PhaseID int

// Phases track different parts of the event (e.g. pools, semis, finals)
type Phase struct {
	ID     PhaseID
	Name   string
	Groups []*Group
}
