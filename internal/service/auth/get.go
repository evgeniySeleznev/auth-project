package auth

import (
	"context"

	"github.com/evgeniySeleznev/auth-project/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	auth, err := s.authRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
