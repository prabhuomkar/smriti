package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaItemsTableName(t *testing.T) {
	mediaItem := MediaItem{}
	assert.Equal(t, MediaItemsTable, mediaItem.TableName())
}
