package storage

import (
	"github.com/groundstorm/cobalt/src/users"
)

// MemoryStorage implements the storage interfaces in memory.  Useful for automated
// testing and rapid iteration during development.  Obviously not suitable for
// production!
type MemoryStorage struct {
	users []memoryStorageUser
}

// NewMemoryStorage creates a new MemoryStorage object
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: []memoryStorageUser{},
	}
}

// AuthenticateUser looks up the user by email and password
func (ms *MemoryStorage) AuthenticateUser(email users.UserEmail, password string) (users.User, error) {
	for _, u := range ms.users {
		if u.Email == email {
			if u.password != password {
				return users.User{}, users.ErrInvalidPassword
			}
			return u.User, nil
		}
	}
	return users.User{}, users.ErrUnknownUser
}

// CreateNewUser creates a new user
func (ms *MemoryStorage) CreateNewUser(user users.User, password string) error {
	for _, u := range ms.users {
		if u.Email == user.Email {
			return users.ErrUserAlreadyExists
		}
	}
	ms.users = append(ms.users, memoryStorageUser{
		User:     user,
		password: password,
	})
	return nil
}
