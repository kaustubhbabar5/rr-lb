package checker

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	wg   sync.WaitGroup
	repo IRepository
}

func New(cache *redis.Client) *Client {
	repo := NewRepository(cache)
	var wg sync.WaitGroup
	return &Client{
		wg:   wg,
		repo: repo,
	}
}

func (c *Client) StartNewHealthCheck(ctx context.Context, url string, period, successThreshold int) {
	c.wg.Add(1)
	go c.Check(ctx, url, 5, 3)
}

func (c *Client) Check(ctx context.Context, url string, period, successThreshold int) {
	defer c.wg.Done()

	ticker := time.NewTicker(time.Second * time.Duration(period))
	successCounter := 0
	healthy := true
	for {
		select {
		case <-ctx.Done():
			log.Info("gracefully exiting health checker for: ", url)
			return
		case <-ticker.C:

			_, err := http.Get(url)
			// _, err := net.DialTimeout("tcp", url+":http", 5)
			if err != nil {
				log.Warning("health check failed for: ", err.Error())
				if !healthy {
					continue
				}
				err = c.repo.MarkReplicaUnhealthy(url)
				if err != nil {
					log.Debug("failed to mark server as unhealthy for: ", url)
					continue
				}
				healthy = false
				log.Debug("marked server as unhealthy for: ", url)

				successCounter = 0
				continue
			}
			log.Debug("health check successful for: ", url)
			if successCounter < 4 {
				successCounter++
			}
			if successCounter == 3 {
				err = c.repo.MarkReplicaHealthy(url)
				if err != nil {
					log.Error("failed to mark server as healthy for: ", url)
					continue
				}
				healthy = true
				log.Info("marked server as healthy for: ", url)

			}
		}
	}

}

func (c *Client) Done() {
	c.wg.Done()
}
