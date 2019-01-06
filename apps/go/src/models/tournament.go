package models

// TournamentID represents a globally unique identifier for this tournament
type TournamentID int

// Tournament tracks big events like Evo.  Many events, many players, etc.
type Tournament struct {
	ID           TournamentID
	Name         string
	Slug         string
	Events       map[EventID]*Event
	Phases       map[PhaseID]*Phase
	Groups       map[GroupID]*Group
	Sets         map[SetID]*Set
	Slots        map[SlotID]*Slot
	Participants map[ParticipantID]*Participant
	Entrants     map[EntrantID]*Entrant
}
