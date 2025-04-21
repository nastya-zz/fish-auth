package converter

import (
	"auth/internal/model"
	desc "auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
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

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	role, _ := strconv.Atoi(user.Role)
	return &desc.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      desc.Role(role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUpdateUserFromDesc(user *desc.UpdateUserInfo) *model.UpdateUser {

	email := ""
	if user.Email != nil {
		email = user.Email.GetValue()
	}

	name := ""
	if user.Name != nil {
		name = user.Name.GetValue()
	}

	return &model.UpdateUser{
		Email: email,
		Name:  name,
		ID:    user.Id,
	}
}
