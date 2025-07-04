package converter

import (
	"auth/internal/model"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
)

func RoleFromDesc(number int) string {
	switch number {
	case 0:
		return "ADMIN"
	case 1:
		return "USER"

	}

	return "USER"
}

func ToCreateUserFromDesc(createUser *desc.UserInfo) *model.CreateUser {
	return &model.CreateUser{
		Email:    createUser.Email,
		Password: createUser.Password,
		Name:     createUser.Name,
		Role:     RoleFromDesc(int(createUser.Role)),
	}
}

func ToUpdateUserFromDesc(id string, updateUser *desc.UserInfo) *model.UpdateUser {
	return &model.UpdateUser{
		ID:         id,
		Email:      updateUser.Email,
		Password:   updateUser.Password,
		Name:       updateUser.Name,
		IsVerified: updateUser.IsVerified,
	}
}
