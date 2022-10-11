package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	_, err := Init()
	assert.Nil(t, err)
}
