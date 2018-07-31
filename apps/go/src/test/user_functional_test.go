package test

import (
	"testing"

	"github.com/groundstorm/cobalt/apps/go/src/models"
	"github.com/groundstorm/cobalt/apps/go/src/storage"
)

func TestNewUser(t *testing.T) {
	ms := storage.NewMemoryStorage()
	UserStorageBasicTests(ms, t)
}

func UserStorageBasicTests(us storage.Storage, t *testing.T) {
	user := models.User{
		FirstName: "first",
		LastName:  "last",
		Email:     "someone@somewhere.com",
	}
	password := "p455w0rd!"

	// try to create the user.
	_, err := us.CreateNewUser(user, password)
	if err != nil {
		t.Errorf("failed to create new user: %v", err)
	}

	// try to create it again.  should fail.
	_, err = us.CreateNewUser(user, password)
	if err != storage.ErrUserAlreadyExists {
		t.Errorf("expected storage.ErrUserAlreadyExists. got: %v", err)
	}

	// see if we can authenticate
	_, err = us.AuthenticateUser(user.Email, password)
	if err != nil {
		t.Errorf("error authenticating new user: %v", err)
	}

	// make sure authentication fails
	_, err = us.AuthenticateUser(user.Email, password+"!")
	if err != storage.ErrInvalidPassword {
		t.Errorf("expected storage.ErrInvalidPassword.  got: %v", err)
	}
}
