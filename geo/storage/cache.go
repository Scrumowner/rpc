package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const Expires = time.Hour

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(db *redis.Client) *RedisCache {
	return &RedisCache{client: db}
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return []byte{}, redis.Nil
	} else if err != nil {
		return []byte{}, fmt.Errorf("Error")
	}
	return val, nil
}
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expires time.Duration) error {
	byte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, key, byte, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}
