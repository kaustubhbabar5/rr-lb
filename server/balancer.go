package server

import (
	"github.com/kaustubhbabar5/rr-lb/balancer"
	"github.com/kaustubhbabar5/rr-lb/stratergy/robin"
)

func (s *HTTPServer) RegisterBalancer() {
	robinRepo := robin.NewRepository(s.cache)
	robinService := robin.NewService(robinRepo)

	repo := balancer.NewRepository(s.cache)
	service := balancer.NewService(repo)

	handler := balancer.NewHandler(service, robinService)
	s.router.HandleFunc("/health", handler.Health).Methods("GET")
	s.router.HandleFunc("/url/register", handler.Register).Methods("POST")
	s.router.HandleFunc("/proxy/{rest:.*}", handler.Proxy)

}
