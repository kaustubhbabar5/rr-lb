package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/kaustubhbabar5/rr-lb/adapters/cache"
)

type HTTPServer struct {
	router *mux.Router
	cache  *redis.Client
}

func New() (*http.Server, error) {

	router := mux.NewRouter()

	//TODO remove magic string
	cache, err := cache.New("0.0.0.0:6379", "")
	if err != nil {
		return nil, fmt.Errorf("redis.NewRedisCache: %w", err)
	}

	server := HTTPServer{
		router: router,
		cache:  cache,
	}

	server.RegisterBalancer()

	return &http.Server{
		// TODO remove magic string
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}, nil
}
