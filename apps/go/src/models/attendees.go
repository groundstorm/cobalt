package models

// Tournament tracks big events like Evo.  Many events, many players, etc.
type Attendees struct {
	Participants []Participant
}

func NewAttendees() *Attendees {
	return &Attendees{
		Participants: []Participant{},
	}
}
