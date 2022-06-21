package balancer

import "fmt"

type IService interface {
	AddNode(url string) error
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) IService {
	return &service{repo}
}

func (s *service) AddNode(url string) error {
	_, err := s.repo.AddReplica(url)
	if err != nil {
		return fmt.Errorf("s.repo.AddReplica: %w", err)
	}
	return nil
}
