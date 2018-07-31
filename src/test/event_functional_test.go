package test

import (
	"fmt"
	"testing"

	"github.com/groundstorm/cobalt/src/event"
	"github.com/groundstorm/cobalt/src/storage"
	"github.com/groundstorm/cobalt/src/users"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	ms := storage.NewMemoryStorage()

	// create a new event
	e := event.Event{
		Slug: "test_event",
	}
	err := ms.CreateEvent(e)
	assert.Nil(t, err)

	// Add a ton of users.
	for i := 0; i < 10; i++ {
		user := users.User{
			FirstName: fmt.Sprintf("First Name %d", i),
			LastName:  fmt.Sprintf("Laste Name %d", i),
			Email:     users.Email(fmt.Sprintf("first%d@last.com", i)),
		}
		err = ms.CreateNewUser(user, "827135871546")
		assert.Nil(t, err)
	}
}
