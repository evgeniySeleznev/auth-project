package converter

import (
	"github.com/evgeniySeleznev/auth-project/internal/model"
	desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
)

func ToDescFromModel(user *model.User) *desc.User {
	return &desc.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     desc.Role(user.Role), // конвертируем наш Role в protobuf Role
	}
}

func ToModelFromDesc(info *desc.CreateRequest) *model.User {
	return &model.User{
		Name:     info.Name,
		Email:    info.Email,
		Password: info.Password,
		Role:     model.Role(info.Role),
	}
}
