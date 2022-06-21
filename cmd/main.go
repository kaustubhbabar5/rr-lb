package main

import (
	"github.com/kaustubhbabar5/rr-lb/server"
)

func main() {
	server, err := server.New()
	if err != nil {
		//TODO handle error here
		panic(err)
	}
	err = server.ListenAndServe()
	if err != nil {
		//TODO handle error properly
		panic(err)
	}
}
