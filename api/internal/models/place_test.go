package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlacesTableName(t *testing.T) {
	place := Place{}
	assert.Equal(t, PlaceTable, place.TableName())
}
