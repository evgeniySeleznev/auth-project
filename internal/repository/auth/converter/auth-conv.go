package converter

import (
	"github.com/evgeniySeleznev/auth-project/internal/repository/auth/model"
	desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
)

func ToUserFromModel(user *model.User) *desc.User {
	return &desc.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     desc.Role(user.Role), // конвертируем наш Role в protobuf Role
	}
}
