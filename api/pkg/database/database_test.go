package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestGetLogLevel(t *testing.T) {
	assert.Equal(t, logger.Error, getLogLevel("ERROR"))
	assert.Equal(t, logger.Warn, getLogLevel("WARN"))
	assert.Equal(t, logger.Info, getLogLevel("INFO"))
	assert.Equal(t, logger.Silent, getLogLevel(""))
}

func TestInit(t *testing.T) {
	db, err := Init("WARNING", "host", 1000, "username", "password", "name")
	assert.Nil(t, db)
	assert.Error(t, err)
}
