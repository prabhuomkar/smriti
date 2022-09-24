package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaItemNewID(t *testing.T) {
	mediaitem := MediaItem{}
	assert.Empty(t, mediaitem.ID)
	mediaitem.NewID()
	assert.NotEmpty(t, mediaitem.ID)
}
