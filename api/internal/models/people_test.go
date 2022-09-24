package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeopleNewID(t *testing.T) {
	people := People{}
	assert.Empty(t, people.ID)
	people.NewID()
	assert.NotEmpty(t, people.ID)
}
