package cache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cache := Init("inmemory", "", 0, "")
	assert.Equal(t, reflect.TypeOf(&InMemoryCache{}), reflect.TypeOf(cache))
	cache = Init("redis", "username", 0, "password")
	assert.Equal(t, reflect.TypeOf(&RedisCache{}), reflect.TypeOf(cache))
}
