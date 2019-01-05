package models

// EventID represents a globally unique identifier for this event
type EventID int

// Event tracks big events like Evo.  Many tournaments, many players, etc.
type Event struct {
	ID     EventID
	Slug   string
	Phases []*Phase
}
