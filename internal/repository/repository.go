package repository

import (
	"context"
	desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthRepository interface {
	Create(ctx context.Context, info *desc.User) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
	Update(ctx context.Context, req *desc.User) (*emptypb.Empty, error)
	Delete(ctx context.Context, req *desc.User) (*emptypb.Empty, error)
}
