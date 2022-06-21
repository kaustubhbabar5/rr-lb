package robin

import "fmt"

type IService interface {
	GetServer() (string, error)
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) IService {
	return &service{repo}
}

func (s *service) GetServer() (string, error) {
	url, err := s.repo.GetServer()
	if err != nil {
		return "", fmt.Errorf("s.repo.GetServer: %w", err)
	}
	return url, nil
}
