package cache

import (
	"fmt"
	"math"
	"time"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

// Provider ...
type Provider interface {
	Ping() error
	SetWithExpire(key string, value any, expiration time.Duration) error
	Get(key string) (any, error)
	Remove(key string) error
}

// Init ...
func Init(cacheType, host string, port int, password string) Provider { //nolint: ireturn
	switch cacheType {
	case "redis":
		addr := fmt.Sprintf("%s:%d", host, port)

		return &RedisCache{
			Connection: &redisClient{
				client: redis.NewClient(&redis.Options{
					Addr:     addr,
					Password: password,
				}),
			},
		}
	default:
		return &InMemoryCache{Connection: gcache.New(math.MaxInt).LRU().Build()}
	}
}
