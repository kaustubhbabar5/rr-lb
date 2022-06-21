package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

//NewRedisCache return a new redis cache
func New(address, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("rdb.Ping: %w", err)
	}

	return rdb, nil
}
