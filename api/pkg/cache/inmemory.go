package cache

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

// InMemoryCache ...
type InMemoryCache struct {
	Connection gcache.Cache
}

func (imc *InMemoryCache) SetWithExpire(key string, value interface{}, expiration time.Duration) error {
	return imc.Connection.SetWithExpire(key, value, expiration)
}

func (imc *InMemoryCache) Get(key string) (interface{}, error) {
	return imc.Connection.Get(key)
}

func (imc *InMemoryCache) Remove(key string) error {
	result := imc.Connection.Remove(key)
	if !result {
		return fmt.Errorf("error removing from cache for key: %+v", key)
	}
	return nil
}
