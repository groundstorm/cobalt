package test

import (
	"fmt"
	"testing"

	"github.com/groundstorm/cobalt/src/models"
	"github.com/groundstorm/cobalt/src/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	ms := storage.NewMemoryStorage()

	// create a new event
	e := models.Event{
		Slug: "test_event",
	}
	eventID, err := ms.CreateEvent(e)
	assert.Nil(t, err)

	// Add a ton of users.
	for i := 0; i < 10; i++ {
		user := models.User{
			FirstName: fmt.Sprintf("First Name %d", i),
			LastName:  fmt.Sprintf("Last Name %d", i),
			Email:     models.Email(fmt.Sprintf("first%d@last.com", i)),
		}
		userID, err := ms.CreateNewUser(user, "827135871546")
		assert.Nil(t, err)
		err = ms.AddUserToEvent(eventID, userID)
		assert.Nil(t, err)
	}
}
