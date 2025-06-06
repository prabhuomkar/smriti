package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	db, err := Init("host", 1000, "username", "password", "name", time.Second)
	assert.Nil(t, db)
	assert.Error(t, err)
}
