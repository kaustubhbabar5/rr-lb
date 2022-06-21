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
	wg         sync.WaitGroup
	repo       IRepository
	httpClient *http.Client
}

type server struct {
	url            string
	healthy        bool
	successCounter int
}

func New(cache *redis.Client, httpClient *http.Client) *Client {
	repo := NewRepository(cache)
	var wg sync.WaitGroup
	return &Client{
		wg:         wg,
		repo:       repo,
		httpClient: httpClient,
	}
}

func (c *Client) StartNewHealthCheck(ctx context.Context, url string, period, successThreshold int) {
	c.wg.Add(1)
	go c.Check(ctx, url, 5, 3)
}

func (c *Client) Check(ctx context.Context, url string, period, successThreshold int) {
	defer c.wg.Done()
	server := server{
		url:            url,
		healthy:        false,
		successCounter: 0,
	}
	ticker := time.NewTicker(time.Second * time.Duration(period))
	successCounter := 0
	for {
		select {
		case <-ctx.Done():
			log.Info("gracefully exiting health checker for: ", url)
			return

		case <-ticker.C:
			ok := c.Ping(server.url)
			if !ok {
				log.Warning("health check failed for: ", server.url)
				if !server.healthy {
					continue
				}

				err := c.repo.MarkReplicaUnhealthy(server.url)
				if err != nil {
					log.Error("failed to mark server as unhealthy for: ", server.url)
					continue
				}

				server.healthy = false
				log.Info("marked server as unhealthy for: ", server.url)
				successCounter = 0
				continue
			}
			log.Debug("health check successful for: ", server.url)
			if successCounter < 4 {
				successCounter++
			}

			if successCounter == 3 {
				err := c.repo.MarkReplicaHealthy(server.url)
				if err != nil {
					log.Error("failed to mark server as healthy for: ", url)
					continue
				}

				server.healthy = true
				log.Info("marked server as healthy for: ", url)
			}
		}
	}

}

func (c *Client) Ping(url string) bool {
	res, err := c.httpClient.Get(url)
	if err != nil {
		log.Debug("failed to ping: ", url, ",err:", err.Error())
		return false
	}
	if res.StatusCode != 200 {
		log.Debug("failed to ping: ", url, ",statusCode:", res.StatusCode)
		return false
	}
	return true

}

func (c *Client) Done() {
	c.wg.Done()
}
