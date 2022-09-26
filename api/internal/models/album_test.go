package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlbumsTableName(t *testing.T) {
	album := Album{}
	assert.Equal(t, AlbumsTable, album.TableName())
}
