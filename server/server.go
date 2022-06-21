package server

import (
	"net/http"
	"os"
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

	server := HTTPServer{
		router:        router,
		cache:         cache,
		checkerClient: checkerClient,
	}

	server.RegisterBalancer()

	return &http.Server{
		Addr:         os.Getenv("HOST") + ":" + os.Getenv("PORT"),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}, nil
}
