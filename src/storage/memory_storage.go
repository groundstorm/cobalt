package storage

import (
	"strconv"

	"github.com/groundstorm/cobalt/src/models"
)

// MemoryStorage implements the storage interfaces in memory.  Useful for automated
// testing and rapid iteration during development.  Obviously not suitable for
// production!
type MemoryStorage struct {
	users      []memoryStorageUser
	nextUserID int
}

// NewMemoryStorage creates a new MemoryStorage object
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: []memoryStorageUser{},
	}
}

// AuthenticateUser looks up the user by email and password
func (ms *MemoryStorage) AuthenticateUser(email models.Email, password string) (models.User, error) {
	for _, u := range ms.users {
		if u.Email == email {
			if u.password != password {
				return models.User{}, ErrInvalidPassword
			}
			return u.User, nil
		}
	}
	return models.User{}, ErrUnknownUser
}

// CreateNewUser creates a new user
func (ms *MemoryStorage) CreateNewUser(user models.User, password string) (models.UserID, error) {
	for _, u := range ms.users {
		if u.Email == user.Email {
			return "", ErrUserAlreadyExists
		}
	}

	ms.nextUserID++
	user.ID = models.UserID(strconv.Itoa(ms.nextUserID))
	ms.users = append(ms.users, memoryStorageUser{
		User:     user,
		password: password,
	})
	return user.ID, nil
}

func (ms *MemoryStorage) LoadEvent(id models.EventID) (models.Event, error) {
	return models.Event{}, nil
}

func (ms *MemoryStorage) CreateEvent(e models.Event) (models.EventID, error) {
	var id models.EventID
	return id, nil
}

func (ms *MemoryStorage) AddUserToEvent(eventID models.EventID, userID models.UserID) error {
	return nil
}
