package converter

import "auth/internal/model"
import modelRepo "auth/internal/repository/auth/model"

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
