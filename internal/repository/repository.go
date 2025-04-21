package repository

import (
	"auth/internal/model"
	"context"
)

type AuthRepository interface {
	Get(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.CreateUser) (int64, error)
	Update(ctx context.Context, updateUser *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email string) (*model.User, error)
}
