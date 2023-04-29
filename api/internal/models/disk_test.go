package models

import (
	"api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDisk(t *testing.T) {
	disk := GetDisk(&config.Config{StorageDiskRoot: "/tmp"})
	assert.NotNil(t, disk)

	disk = GetDisk(&config.Config{StorageDiskRoot: "invalid"})
	assert.Nil(t, disk)
}
