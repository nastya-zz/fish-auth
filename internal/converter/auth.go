package converter

import (
	"auth/internal/model"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
)

func role(number int) string {
	switch number {
	case 0:
		return "ADMIN"
	case 1:
		return "USER"
	}

	return "USER"
}

func ToCreateUserFromDesc(createUser *desc.CreateUserInfo) *model.CreateUser {
	return &model.CreateUser{
		Email:    createUser.Email,
		Password: createUser.Password,
		Name:     createUser.Name,
		Role:     role(int(createUser.Role)),
	}
}
