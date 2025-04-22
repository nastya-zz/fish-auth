package repository

import (
	"auth/internal/model"
	"context"
)

type AuthRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user *model.CreateUser) (string, error)
	Update(ctx context.Context, updateUser *model.UpdateUser) error
	Delete(ctx context.Context, id string) error
	Block(ctx context.Context, id string) error
	Login(ctx context.Context, email string) (*model.User, error)
}
