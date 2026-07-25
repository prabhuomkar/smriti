package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	health := GetHealth()
	assert.Equal(t, &Health{Status: StatusUp}, health)
}
