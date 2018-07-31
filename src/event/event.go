package event

// ID represents a globally unique identifier for this event
type ID string

// Event tracks big events like Evo.  Many tournaments, many players, etc.
type Event struct {
	ID   ID
	Slug string
}
