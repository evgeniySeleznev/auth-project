package service

import (
	"context"

	"github.com/evgeniySeleznev/auth-project/internal/model"
)

type AuthService interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}
