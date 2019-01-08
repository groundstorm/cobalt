package config

// CobaltConfig contains global cobalt configuration info
type CobaltConfig struct {
	Tournaments []TournamentConfig
}

func NewCobaltConfig() CobaltConfig {
	return CobaltConfig{
		Tournaments: []TournamentConfig{},
	}
}
