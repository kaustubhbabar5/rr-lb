package server

import (
	"github.com/kaustubhbabar5/rr-lb/balancer"
)

func (s *HTTPServer) RegisterBalancer() {
	repo := balancer.NewRepository(s.cache)
	service := balancer.NewService(repo, s.checkerClient)

	handler := balancer.NewHandler(service)
	s.router.HandleFunc("/health", handler.Health).Methods("GET")
	s.router.HandleFunc("/url/register", handler.Register).Methods("POST")
	s.router.HandleFunc("/proxy/{rest:.*}", handler.Proxy)

}
