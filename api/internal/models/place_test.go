package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaceNewID(t *testing.T) {
	place := Place{}
	assert.Empty(t, place.ID)
	place.NewID()
	assert.NotEmpty(t, place.ID)
}
