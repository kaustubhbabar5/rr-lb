package robin

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/rr-lb/pkg/constants"
)

type IRepository interface {
	GetServer() (string, error)
}

type repository struct {
	cache *redis.Client
}

func NewRepository(cache *redis.Client) *repository {
	return &repository{cache}
}

//
func (r *repository) GetServer() (string, error) {
	ctx := context.Background()

	res := r.cache.LMove(ctx, constants.HEALTHY_SERVERS, constants.HEALTHY_SERVERS, "LEFT", "RIGHT")
	if res.Err() != nil {
		return "", fmt.Errorf("tx.LMove: %w", res.Err())
	}

	url := res.Val()

	return url, nil
}
