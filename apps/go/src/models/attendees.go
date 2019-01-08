package models

// Tournament tracks big events like Evo.  Many events, many players, etc.
type Attendees struct {
	Registrations map[ParticipantID]Registration
	Events        map[EventID]Event
}
