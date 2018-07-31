package models

// ID represents a globally unique identifier for this tournament
type TournamentID string

// Tournament tracks a single tournament at an event
type Tournament struct {
	ID   TournamentID
	Slug string
}
