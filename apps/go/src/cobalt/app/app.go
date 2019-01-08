package app

import (
	"encoding/json"
	"fmt"

	"github.com/groundstorm/cobalt/apps/go/src/cobalt/app/config"
	"github.com/groundstorm/cobalt/apps/go/src/models"

	"github.com/boltdb/bolt"
	logging "github.com/op/go-logging"
)

var (
	log  = logging.MustGetLogger("app")
	slog = logging.MustGetLogger("app.storage")
)

type App struct {
	db     *bolt.DB
	config config.CobaltConfig
}

func New() (*App, error) {
	log.Debugf("opening bolt...")
	db, err := bolt.Open("cobalt.db", 0600, nil)
	log.Debugf("done!\n")
	if err != nil {
		return nil, err
	}
	a := &App{
		db:     db,
		config: config.NewCobaltConfig(),
	}

	err = a.loadConfig()
	if err != nil {
		log.Infof("error loading config (%s).  creating new config.", err)
		err := a.ModifyConfig(func(c *config.CobaltConfig) error {
			*c = config.NewCobaltConfig()
			return nil
		})
		if err != nil {
			log.Infof("failed to save new config (%s).", err)
			return nil, err
		}
	}
	return a, nil
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

func (a *App) ModifyConfig(fn func(*config.CobaltConfig) error) error {
	err := fn(&a.config)
	if err != nil {
		return err
	}
	return a.db.Update(func(tx *bolt.Tx) error {
		b, err := makeBucketForConfig(tx)
		if err != nil {
			return err
		}
		value, err := json.Marshal(&a.config)
		if err != nil {
			return fmt.Errorf("failed to marshall config: %s", err)
		}
		b.Put([]byte("config"), value)
		return nil
	})
}

func (a *App) GetConfig() config.CobaltConfig {
	return a.config
}

func (a *App) loadConfig() error {
	return a.db.View(func(tx *bolt.Tx) error {
		b := getBucketForConfig(tx)
		if b == nil {
			return fmt.Errorf("config bucket not found")
		}
		value := b.Get([]byte("config"))
		return json.Unmarshal(value, &a.config)
	})
}
