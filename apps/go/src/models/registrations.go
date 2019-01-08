package models

// Tournament tracks big events like Evo.  Many events, many players, etc.
type Registrations struct {
	Registrations map[ParticipantID]Registration
	Events        map[EventID]Event
}
