package models

// TournamentID represents a globally unique identifier for this tournament
type TournamentID string

// Tournament tracks big events like Evo.  Many events, many players, etc.
type Tournament struct {
	ID   TournamentID
	Slug string
}
