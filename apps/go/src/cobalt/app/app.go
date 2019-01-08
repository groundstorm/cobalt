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
func (a *App) StoreRegs(slug string, regs *models.Registrations) error {
	return a.db.Update(func(tx *bolt.Tx) error {
		ab, err := makeBucketForRegistrations(tx, slug)
		if err != nil {
			log.Errorf(`failed getting bucket for "%s": %s`, slug, err)
			return err
		}

		// remove all existing registrations
		ab.ForEach(func(k, v []byte) error {
			ab.Delete(k)
			return nil
		})
		// add the new ones
		for _, r := range regs.Registrations {
			value, err := json.Marshal(r)
			if err != nil {
				return fmt.Errorf("failed to marshall participant %d: %s", r.Participant.ID, err)
			}
			ab.Put(r.Participant.Key(), value)
		}
		return nil
	})
}

func (a *App) GetRegs(slug string) (*models.Registrations, error) {
	regs := &models.Registrations{
		Registrations: map[models.ParticipantID]models.Registration{},
		Events:        map[models.EventID]models.Event{},
	}
	err := a.db.View(func(tx *bolt.Tx) error {
		ab := getBucketForRegistrations(tx, slug)
		if ab == nil {
			return fmt.Errorf("registrations for %s have not been fetched.", slug)
		}
		return ab.ForEach(func(k, v []byte) error {
			var r models.Registration
			err := json.Unmarshal(v, &r)
			if err != nil {
				return err
			}
			regs.Registrations[r.Participant.ID] = r
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return regs, nil
}
