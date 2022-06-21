package checker

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/rr-lb/pkg/constants"
)

type IRepository interface {
	MarkReplicaUnhealthy(url string) error
	MarkReplicaHealthy(url string) error
}

type repository struct {
	cache *redis.Client
}

func NewRepository(cache *redis.Client) *repository {
	return &repository{cache}
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
