package adapter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is an adapter that wraps Redis client to provide caching functionality
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache instance with the specified Redis connection options
func NewRedisCache(ctx context.Context, addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{
		client: client,
	}
}

// Set stores a value in Redis with an optional expiration duration
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value from Redis. Returns an error if the key does not exist
func (r *RedisCache) Get(ctx context.Context, key string) (*string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &val, nil
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis client connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}
