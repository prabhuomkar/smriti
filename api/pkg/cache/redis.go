package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type (
	// RedisClient ...
	RedisClient interface {
		Get(context.Context, string) (string, error)
		Set(context.Context, string, interface{}, time.Duration) error
		Del(context.Context, string) error
	}

	// RedisCache ...
	RedisCache struct {
		Connection RedisClient
	}

	redisClient struct {
		client *redis.Client
	}
)

func (rc *RedisCache) SetWithExpire(key string, value interface{}, expiration time.Duration) error {
	return rc.Connection.Set(context.TODO(), key, value, expiration)
}

func (rc *RedisCache) Get(key string) (interface{}, error) {
	return rc.Connection.Get(context.TODO(), key)
}

func (rc *RedisCache) Remove(key string) error {
	return rc.Connection.Del(context.TODO(), key)
}

func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisClient) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
