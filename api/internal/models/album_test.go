package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlbumNewID(t *testing.T) {
	album := Album{}
	assert.Empty(t, album.ID)
	album.NewID()
	assert.NotEmpty(t, album.ID)
}
