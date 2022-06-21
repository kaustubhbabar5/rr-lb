package balancer

import (
	"context"
	"fmt"

	"github.com/kaustubhbabar5/rr-lb/checker"
)

type IService interface {
	AddServer(url string) error
	GetServer() (string, error)
}

type service struct {
	repo          IRepository
	checkerClient *checker.Client
}

func NewService(repo IRepository, checkerClient *checker.Client) IService {
	return &service{repo, checkerClient}
}

func (s *service) AddServer(url string) error {
	_, err := s.repo.AddServer(url)
	if err != nil {
		return fmt.Errorf("s.repo.AddReplica: %w", err)
	}

	s.checkerClient.StartNewHealthCheck(context.Background(), url, 5, 3)

	return nil
}

func (s *service) GetServer() (string, error) {
	url, err := s.repo.GetServer()
	if err != nil {
		return "", fmt.Errorf("s.repo.GetServer: %w", err)
	}
	return url, nil
}
