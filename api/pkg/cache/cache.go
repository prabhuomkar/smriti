package cache

import (
	"api/config"
	"fmt"
	"math"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
)

// Provider ...
type Provider interface {
	SetWithExpire(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Remove(key string) error
}

// Init ...
func Init(config *config.Config) Provider { //nolint: ireturn
	switch config.Cache.Type {
	case "redis":
		return &RedisCache{
			Connection: &redisClient{client: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", config.Cache.Host, config.Cache.Port),
				Password: config.Cache.Password,
			})},
		}
	default:
		return &InMemoryCache{
			Connection: gcache.New(math.MaxInt).
				LRU().
				Build(),
		}
	}
}
