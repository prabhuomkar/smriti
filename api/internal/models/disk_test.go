package models

import (
	"api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDisk(t *testing.T) {
	disk := GetDisk(&config.Config{Storage: config.Storage{DiskRoot: "/tmp"}})
	assert.NotNil(t, disk)

	disk = GetDisk(&config.Config{Storage: config.Storage{DiskRoot: "invalid"}})
	assert.Nil(t, disk)
}
