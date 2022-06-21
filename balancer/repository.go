package balancer

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/rr-lb/pkg/constants"
	cerrors "github.com/kaustubhbabar5/rr-lb/pkg/errors"
)

type IRepository interface {
	AddServer(url string) (int64, error)
	GetServer() (string, error)
}

type repository struct {
	cache *redis.Client
}

func NewRepository(cache *redis.Client) *repository {
	return &repository{cache}
}

//checks if key exists in list
func (r *repository) Exists(key, value string) (bool, error) {
	res := r.cache.LPos(context.Background(), key, value, redis.LPosArgs{})
	if res.Err() != nil {
		if res.Err().Error() == "redis: nil" {
			//TODO handle properly
			return false, nil
		}
		return false, fmt.Errorf("r.cache.LPos: %w", res.Err())
	}
	return true, nil
}

func (r *repository) AddServer(url string) (int64, error) {
	ctx := context.Background()

	exists, err := r.Exists(constants.UNHEALTHY_SERVERS, url)
	if err != nil {
		return 0, fmt.Errorf("r.Exists: %w", err)
	}

	if exists {
		return 0, errors.New("server already registered")
	}

	addRes := r.cache.LPush(ctx, constants.UNHEALTHY_SERVERS, url)
	if addRes.Err() != nil {
		return 0, fmt.Errorf("tx.LPush: %w", addRes.Err())
	}

	index := addRes.Val()

	return index, nil
}

func (r *repository) GetServer() (string, error) {
	ctx := context.Background()

	count, err := r.GetHealthyServerCount()
	if err != nil {
		return "", fmt.Errorf("r.GetHealthyServerCount: %w", err)
	}

	if count == int64(0) {
		return "", cerrors.NewNotFound("healthy_server", "0 healthy servers found")
	}

	res := r.cache.LMove(ctx, constants.HEALTHY_SERVERS, constants.HEALTHY_SERVERS, "LEFT", "RIGHT")
	if res.Err() != nil {
		return "", fmt.Errorf("r.cache.LMove: %w", res.Err())
	}

	url := res.Val()

	return url, nil
}

func (r *repository) GetHealthyServerCount() (int64, error) {
	ctx := context.Background()
	res := r.cache.LLen(ctx, constants.HEALTHY_SERVERS)
	if res.Err() != nil {
		return 0, fmt.Errorf("r.cache.LLen %w", res.Err())
	}
	return res.Val(), nil
}
