package storage

import (
	"testing"

	"github.com/groundstorm/cobalt/src/users"
)

func TestNewUser(t *testing.T) {
	ms := NewMemoryStorage()
	UserStorageBasicTests(ms, t)
}

func UserStorageBasicTests(us users.UserStorage, t *testing.T) {
	user := users.User{
		FirstName: "first",
		LastName:  "last",
		Email:     "someone@somewhere.com",
	}
	password := "p455w0rd!"

	// try to create the user.
	err := us.CreateNewUser(user, password)
	if err != nil {
		t.Errorf("failed to create new user: %v", err)
	}

	// try to create it again.  should fail.
	err = us.CreateNewUser(user, password)
	if err != users.ErrUserAlreadyExists {
		t.Errorf("expected users.ErrUserAlreadyExists. got: %v", err)
	}

	// see if we can authenticate
	_, err = us.AuthenticateUser(user.Email, password)
	if err != nil {
		t.Errorf("error authenticating new user: %v", err)
	}

	// make sure authentication fails
	_, err = us.AuthenticateUser(user.Email, password+"!")
	if err != users.ErrInvalidPassword {
		t.Errorf("expected users.ErrInvalidPassword.  got: %v", err)
	}
}
