package auth

import (
	"context"

	"github.com/evgeniySeleznev/auth-project/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.User) (int64, error) {
	id, err := s.authRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}

	return id, nil
}
