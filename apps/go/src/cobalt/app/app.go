package app

import (
	"encoding/json"
	"fmt"

	"github.com/groundstorm/cobalt/apps/go/src/models"

	"github.com/boltdb/bolt"
	logging "github.com/op/go-logging"
)

var (
	log  = logging.MustGetLogger("app")
	slog = logging.MustGetLogger("app.storage")
)

type App struct {
	db *bolt.DB
}

func New() (*App, error) {
	fmt.Printf("opening bolt...")
	db, err := bolt.Open("cobalt.db", 0600, nil)
	fmt.Printf("done!\n")
	if err != nil {
		return nil, err
	}
	return &App{
		db: db,
	}, nil
}

func (a *App) Close() {
	a.db.Close()
}

// ImportAttendees will import the attendee list for a tournament into
// the db.
func (a *App) StoreRegs(slug string, attendees *models.Attendees) error {
	return a.db.Update(func(tx *bolt.Tx) error {
		ab, err := makeBucketForAttendees(tx, slug)
		if err != nil {
			log.Errorf(`failed getting bucket for "%s": %s`, slug, err)
			return err
		}

		// remove all existing attendees
		ab.ForEach(func(k, v []byte) error {
			ab.Delete(k)
			return nil
		})
		// add the new ones
		for _, p := range attendees.Participants {
			value, err := json.Marshal(p)
			if err != nil {
				return fmt.Errorf("failed to marshall participant %d: %s", p.ID, err)
			}
			ab.Put(p.Key(), value)
		}
		return nil
	})
}

func (a *App) LoadRegs(slug string) (*models.Attendees, error) {
	regs := models.NewAttendees()
	err := a.db.View(func(tx *bolt.Tx) error {
		ab := getBucketForAttendees(tx, slug)
		if ab == nil {
			return fmt.Errorf("registrations for %s have not been fetched.", slug)
		}
		return ab.ForEach(func(k, v []byte) error {
			var p models.Participant
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			regs.Participants = append(regs.Participants, p)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return regs, nil
}
