package models

import (
	"strconv"
)

// ParticipantID represents a globally unique identifier for this user
type ParticipantID int

// The Participant struct represents any user of the system.  Players, spectators,
// TOs, Judges, etc.  All must have a user account to do anything
type Participant struct {
	ID     ParticipantID `json:"id"`
	Player Player        `json:"player"`
}

func (p *Participant) Key() []byte {
	return []byte(strconv.Itoa(int(p.ID)))
}
