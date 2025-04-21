package service

import (
	"auth/internal/model"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	GetAccessToken(ctx context.Context, claims *model.UserClaims) (string, error)
	GetRefreshToken(ctx context.Context, claims *model.UserClaims) (string, error)
}

type UserService interface {
	Create(ctx context.Context, user *model.CreateUser) (int64, error)
	Update(ctx context.Context, updateUser *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
	//BlockUser(ctx context.Context, id int64) error
}
