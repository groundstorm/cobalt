package models

type Registration struct {
	Participant Participant
	Events      map[EventID]bool
}
