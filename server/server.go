package server

import (
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/kaustubhbabar5/rr-lb/checker"
)

type HTTPServer struct {
	router        *mux.Router
	cache         *redis.Client
	checkerClient *checker.Client
}

func New(cache *redis.Client, checkerClient *checker.Client) (*http.Server, error) {

	router := mux.NewRouter()

	//TODO remove magic string

	server := HTTPServer{
		router:        router,
		cache:         cache,
		checkerClient: checkerClient,
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
