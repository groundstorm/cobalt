package test

import (
	"fmt"
	"testing"

	"github.com/groundstorm/cobalt/apps/go/src/smashgg"
	"github.com/stretchr/testify/assert"
)

func TestSmashGGImport(t *testing.T) {
	tourney, err := smashgg.ImportTournament("evo2018")
	fmt.Printf("%v\n", tourney)
	assert.Nil(t, err)
	assert.NotNil(t, tourney)
}
