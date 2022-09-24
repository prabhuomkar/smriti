package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThingNewID(t *testing.T) {
	thing := Thing{}
	assert.Empty(t, thing.ID)
	thing.NewID()
	assert.NotEmpty(t, thing.ID)
}
