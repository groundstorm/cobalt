package event

import (
	"fmt"
)

// The Storage interface is used to load and store user information
type Storage interface {
	LoadEvent(id ID) (Event, error)
	CreateEvent(e Event) error
}

var (
	ErrUnknonwnEvent  = fmt.Errorf("UNKNOWN_EVENT")
	ErrDuplicateEvent = fmt.Errorf("DUPLICATE_EVENT")
)
