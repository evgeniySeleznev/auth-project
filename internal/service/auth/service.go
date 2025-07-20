package auth

import "github.com/evgeniySeleznev/auth-project/internal/repository"

import (
	def "github.com/evgeniySeleznev/auth-project/internal/service"
)

var _ def.AuthService = (*serv)(nil)

type serv struct {
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository) *serv {
	return &serv{
		authRepository: authRepository,
	}
}
