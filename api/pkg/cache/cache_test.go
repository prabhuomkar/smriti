package cache

import (
	"api/config"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cache := Init(&config.Config{Cache: config.Cache{Type: "inmemory"}})
	assert.Equal(t, reflect.TypeOf(&InMemoryCache{}), reflect.TypeOf(cache))
	cache = Init(&config.Config{Cache: config.Cache{Type: "redis"}})
	assert.Equal(t, reflect.TypeOf(&RedisCache{}), reflect.TypeOf(cache))
}
