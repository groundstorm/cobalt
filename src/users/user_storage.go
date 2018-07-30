package users

import (
	"fmt"
)

// The UserStorage interface is used to load and store user information
type UserStorage interface {
	AuthenticateUser(email UserEmail, password string) (User, error)
	CreateNewUser(user User, password string) error
}

var (
	ErrInvalidPassword   = fmt.Errorf("INVALID_PASSWORD")
	ErrUnknownUser       = fmt.Errorf("UNKNOWN_USER")
	ErrUserAlreadyExists = fmt.Errorf("USER_ALREADY_EXISTS")
)
