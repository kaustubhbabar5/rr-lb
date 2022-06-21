package balancer

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/rr-lb/pkg/constants"
)

type IRepository interface {
	AddServer(url string) (int64, error)
	MarkReplicaUnhealthy(url string) error
	MarkReplicaHealthy(url string) error
	GetServer() (string, error)
}

type repository struct {
	cache *redis.Client
}

func NewRepository(cache *redis.Client) *repository {
	return &repository{cache}
}

func (r *repository) AddServer(url string) (int64, error) {
	ctx := context.Background()
	tx := r.cache.TxPipeline()

	// respond back with
	removeRes := tx.LRem(ctx, constants.HEALTHY_SERVERS, 1, url)
	if removeRes.Err() != nil {
		tx.Discard()
		return 0, fmt.Errorf("tx.LRem: %w", removeRes.Err())
	}

	pushRes := tx.LPush(ctx, constants.HEALTHY_SERVERS, url)
	if pushRes.Err() != nil {
		tx.Discard()
		return 0, fmt.Errorf("tx.LPush: %w", pushRes.Err())
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("tx.Exec %w", err)
	}

	index := pushRes.Val()
	return index, nil
}

func (r *repository) MarkReplicaUnhealthy(url string) error {
	ctx := context.Background()
	tx := r.cache.TxPipeline()

	res := tx.LRem(ctx, constants.HEALTHY_SERVERS, 1, url)
	if res.Err() != nil {
		tx.Discard()
		return fmt.Errorf("tx.LRem: %w", res.Err())
	}

	res = tx.RPush(ctx, constants.UNHEALTHY_SERVERS, url)
	if res.Err() != nil {
		tx.Discard()
		return fmt.Errorf("tx.RPush: %w", res.Err())
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf("tx.Exec %w", err)
	}
	return nil
}
func (r *repository) MarkReplicaHealthy(url string) error {
	ctx := context.Background()
	tx := r.cache.TxPipeline()

	res := tx.LRem(ctx, constants.UNHEALTHY_SERVERS, 1, url)
	if res.Err() != nil {
		tx.Discard()
		return fmt.Errorf("tx.LRem: %w", res.Err())
	}

	res = tx.RPush(ctx, constants.HEALTHY_SERVERS, url)
	if res.Err() != nil {
		tx.Discard()
		return fmt.Errorf("tx.RPush: %w", res.Err())
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf("tx.Exec %w", err)
	}
	return nil
}

func (r *repository) GetServer() (string, error) {
	ctx := context.Background()

	res := r.cache.LMove(ctx, constants.HEALTHY_SERVERS, constants.HEALTHY_SERVERS, "LEFT", "RIGHT")
	if res.Err() != nil {
		return "", fmt.Errorf("r.cache.LMove: %w", res.Err())
	}

	url := res.Val()

	return url, nil
}
