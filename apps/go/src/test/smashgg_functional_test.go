package test

import (
	"testing"

	"github.com/groundstorm/cobalt/apps/go/src/smashgg"
	"github.com/stretchr/testify/assert"
)

func TestSmashGGImport(t *testing.T) {
	evt, err := smashgg.ImportEvent("evo2018")
	assert.Nil(t, err)
	assert.NotNil(t, evt)
}
