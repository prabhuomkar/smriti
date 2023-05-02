package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/bluele/gcache"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache(t *testing.T) {
	// error
	errCache := InMemoryCache{
		Connection: gcache.New(1024).
			LRU().
			SerializeFunc(func(a interface{}, b interface{}) (interface{}, error) {
				return nil, errors.New("some error")
			}).
			DeserializeFunc(func(a interface{}, b interface{}) (interface{}, error) {
				return nil, errors.New("some error")
			}).
			Build(),
	}
	assert.Error(t, errCache.SetWithExpire("key", "value", 1*time.Minute))
	val, err := errCache.Get("key")
	assert.Error(t, err)
	assert.Nil(t, val)
	assert.Error(t, errCache.Remove("key"))

	// success
	cache := InMemoryCache{Connection: gcache.New(1024).LRU().Build()}
	assert.NoError(t, cache.SetWithExpire("key", "value", 1*time.Minute))
	val, err = cache.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
	assert.NoError(t, cache.Remove("key"))
}
