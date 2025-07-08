package repository

import (
	"auth/internal/model"
	"context"
)

type AuthRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user *model.CreateUser) (*model.User, error)
	Update(ctx context.Context, updateUser *model.UpdateUser) (*model.UpdateUser, error)
	UpdateRole(ctx context.Context, id, role string) error
	Delete(ctx context.Context, id string) error
	Block(ctx context.Context, id string) error
	Login(ctx context.Context, email string) (*model.User, error)
}

type EventRepository interface {
	GetNewEvent(ctx context.Context, count int) ([]*model.Event, error)
	SaveEvent(ctx context.Context, event *model.Event) error
	SetDone(ctx context.Context, id int) error
}
