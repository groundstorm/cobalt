package storage

import (
	"fmt"

	"github.com/groundstorm/cobalt/apps/go/src/models"
)

// The Storage interface is used to load and store user information
type Storage interface {
	CreateEvent(e models.Event) (models.EventID, error)
	LoadEvent(id models.EventID) (models.Event, error)

	CreateNewUser(user models.User, password string) (models.UserID, error)
	AuthenticateUser(email models.Email, password string) (models.User, error)

	AddUserToEvent(eventID models.EventID, userID models.UserID) error
}

var (
	ErrUnknownEvent      = fmt.Errorf("UNKNOWN_EVENT")
	ErrDuplicateEvent    = fmt.Errorf("DUPLICATE_EVENT")
	ErrInvalidPassword   = fmt.Errorf("INVALID_PASSWORD")
	ErrUnknownUser       = fmt.Errorf("UNKNOWN_USER")
	ErrUserAlreadyExists = fmt.Errorf("USER_ALREADY_EXISTS")
)
