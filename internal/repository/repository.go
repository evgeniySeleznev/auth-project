package repository

import (
	"context"
	//desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
	"github.com/evgeniySeleznev/auth-project/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthRepository interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, req *model.User) (*emptypb.Empty, error)
	Delete(ctx context.Context, req *model.User) (*emptypb.Empty, error)
}
