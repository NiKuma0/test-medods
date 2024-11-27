package services

import "jwt-service/internal/repositories"

type UserService struct {
	repos *repositories.Repositories
}

func NewUserService(repos *repositories.Repositories) *UserService {
	return &UserService{
		repos: repos,
	}
}

func (s *UserService) Get(userId string) (repositories.User, error) {
	return s.repos.User.Get(userId)
}
