package service

import (
	"auth/internal/model"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, error)
	GetAccessToken(ctx context.Context, claims *model.UserClaims) (string, error)
	GetRefreshToken(ctx context.Context, claims *model.UserClaims) (string, error)
	Create(ctx context.Context, user *model.CreateUser) (string, error)
	Check(ctx context.Context, endpoint string) (bool, error)
	BlockUser(ctx context.Context, id string) (string, error)
}
