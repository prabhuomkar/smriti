package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModels(t *testing.T) {
	models := GetModels()
	assert.Len(t, models, 9)
	assert.Equal(t, []interface{}{
		User{},
		Album{},
		Place{},
		Thing{},
		People{},
		MediaItem{},
		MediaitemEmbedding{},
		MediaitemFace{},
		Job{},
	}, models)
}
