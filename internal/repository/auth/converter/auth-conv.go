package converter

import (
	"github.com/evgeniySeleznev/auth-project/internal/model"
	modelRepo "github.com/evgeniySeleznev/auth-project/internal/repository/auth/model"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     model.Role(user.Role), // конвертируем наш Role в protobuf Role
	}
}
