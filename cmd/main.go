package main

import (
	"github.com/kaustubhbabar5/rr-lb/adapters/cache"
	"github.com/kaustubhbabar5/rr-lb/checker"
	"github.com/kaustubhbabar5/rr-lb/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	//
	log.SetLevel(log.InfoLevel)

	cache, err := cache.New("0.0.0.0:6379", "")
	if err != nil {
		panic(err)
	}

	checkerClient := checker.New(cache)

	server, err := server.New(cache, checkerClient)
	if err != nil {
		//TODO handle error here
		panic(err)
	}
	log.Info("starting server")
	err = server.ListenAndServe()
	if err != nil {
		//TODO handle error properly
		panic(err)
	}
	//TODO graceful shutdown on sigterm
	// checkerClient.Done()
}
