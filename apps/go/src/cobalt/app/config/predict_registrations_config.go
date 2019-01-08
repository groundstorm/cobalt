package config

import (
	"github.com/groundstorm/cobalt/apps/go/src/models"
)

// PredictRegistrationsConfig contains info on how we want to predict attendance
type PredictRegistrationsConfig struct {
	CompareTo []models.TournamentID
}
