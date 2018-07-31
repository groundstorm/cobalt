package storage

import (
	"strconv"

	"github.com/groundstorm/cobalt/src/event"
	"github.com/groundstorm/cobalt/src/users"
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
func (ms *MemoryStorage) AuthenticateUser(email users.Email, password string) (users.User, error) {
	for _, u := range ms.users {
		if u.Email == email {
			if u.password != password {
				return users.User{}, ErrInvalidPassword
			}
			return u.User, nil
		}
	}
	return users.User{}, ErrUnknownUser
}

// CreateNewUser creates a new user
func (ms *MemoryStorage) CreateNewUser(user users.User, password string) (users.ID, error) {
	for _, u := range ms.users {
		if u.Email == user.Email {
			return "", ErrUserAlreadyExists
		}
	}

	ms.nextUserID++
	user.ID = users.ID(strconv.Itoa(ms.nextUserID))
	ms.users = append(ms.users, memoryStorageUser{
		User:     user,
		password: password,
	})
	return user.ID, nil
}

func (ms *MemoryStorage) LoadEvent(id event.ID) (event.Event, error) {
	return event.Event{}, nil
}

func (ms *MemoryStorage) CreateEvent(e event.Event) error {
	return nil
}
