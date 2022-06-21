package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kaustubhbabar5/rr-lb/adapters/cache"
	"github.com/kaustubhbabar5/rr-lb/checker"
	"github.com/kaustubhbabar5/rr-lb/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()

	log.SetLevel(log.InfoLevel)

	cache, err := cache.New(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"))
	if err != nil {
		log.Fatalln(err)
	}

	httpClient := &http.Client{
		Timeout: 4 * time.Second,
	}

	checkerClient := checker.New(cache, httpClient)
	//TODO start health checkers on servers that stored in redis on startup

	server, err := server.New(cache, checkerClient)
	if err != nil {
		//TODO handle error here
		log.Fatalln(err)
	}
	log.Info("starting server")
	err = server.ListenAndServe()
	if err != nil {
		//TODO handle error properly
		log.Fatalln(err)
	}
	//TODO graceful shutdown on sigterm
	// checkerClient.Done()
}
