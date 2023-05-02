package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRedisClient struct {
	fail  bool
	value string
}

func (m *mockRedisClient) Get(context.Context, string) (string, error) {
	if m.fail {
		return "", errors.New("some error")
	}
	return m.value, nil
}

func (m *mockRedisClient) Set(context.Context, string, interface{}, time.Duration) error {
	if m.fail {
		return errors.New("some error")
	}
	return nil
}

func (m *mockRedisClient) Del(context.Context, string) error {
	if m.fail {
		return errors.New("some error")
	}
	return nil
}

func TestRedisCache(t *testing.T) {
	// error
	errCache := RedisCache{
		Connection: &mockRedisClient{fail: true},
	}
	assert.Error(t, errCache.SetWithExpire("key", "value", 1*time.Minute))
	val, err := errCache.Get("key")
	assert.Error(t, err)
	assert.Equal(t, "", val)
	assert.Error(t, errCache.Remove("key"))

	// success
	cache := RedisCache{
		Connection: &mockRedisClient{fail: false, value: "value"},
	}
	assert.NoError(t, cache.SetWithExpire("key", "value", 1*time.Minute))
	val, err = cache.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
	assert.NoError(t, cache.Remove("key"))
}
