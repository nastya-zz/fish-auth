package converter

import (
	"auth/internal/model"
	"database/sql"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"time"
)
import modelRepo "auth/internal/repository/auth/model"

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Role:       user.Role,
		Password:   user.Password,
		IsBlocked:  user.IsBlocked,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		LastLogin:  timeConvert(user.LastLogin),
	}
}

func timeConvert(sourceTime sql.NullTime) time.Time {
	if !sourceTime.Valid {
		return time.Time{}
	} else {
		return sourceTime.Time
	}
}

func DescRole(r string) desc.Role {
	switch r {
	case model.RoleAdmin:
		return desc.Role_ADMIN
	case model.RoleUser:
		return desc.Role_USER
	}

	return desc.Role_USER
}
