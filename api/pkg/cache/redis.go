package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ( // RedisClient ...
	RedisClient interface {
		Ping(ctx context.Context) error
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, val any, ttl time.Duration) error
		Del(ctx context.Context, key string) error
	}

	// RedisCache ...
	RedisCache struct {
		Connection RedisClient
	}

	redisClient struct {
		client *redis.Client
	}
)

func (rc *RedisCache) Ping() error {
	return rc.Connection.Ping(context.TODO())
}

func (rc *RedisCache) SetWithExpire(key string, value any, expiration time.Duration) error {
	return rc.Connection.Set(context.TODO(), key, value, expiration)
}

func (rc *RedisCache) Get(key string) (any, error) {
	return rc.Connection.Get(context.TODO(), key)
}

func (rc *RedisCache) Remove(key string) error {
	return rc.Connection.Del(context.TODO(), key)
}

func (c *redisClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisClient) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
