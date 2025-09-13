package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

// InMemoryCache ...
type InMemoryCache struct {
	Connection gcache.Cache
}

var errRemovingFromCache = errors.New("error removing from cache")

func (imc *InMemoryCache) SetWithExpire(key string, value interface{}, expiration time.Duration) error {
	return imc.Connection.SetWithExpire(key, value, expiration)
}

func (imc *InMemoryCache) Get(key string) (interface{}, error) {
	return imc.Connection.Get(key)
}

func (imc *InMemoryCache) Remove(key string) error {
	result := imc.Connection.Remove(key)
	if !result {
		return fmt.Errorf("%w for key: %+v", errRemovingFromCache, key)
	}

	return nil
}
