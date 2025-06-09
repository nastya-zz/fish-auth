package service

import (
	"auth/internal/model"
	"context"
	"time"
)

type AuthService interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	Login(ctx context.Context, login string, password string) (string, error)
	GetAccessToken(ctx context.Context, claims *model.UserClaims) (string, error)
	GetRefreshToken(ctx context.Context, claims *model.UserClaims) (string, error)
	Create(ctx context.Context, user *model.CreateUser) (string, error)
	Check(ctx context.Context, endpoint string) (bool, error)
	BlockUser(ctx context.Context, id string) (string, error)
	Delete(ctx context.Context, id string) error
}

type ProcessEvent interface {
	StartProcessEvents(ctx context.Context, handlePeriod time.Duration)
	SendMessage(ctx context.Context, event []*model.Event)
}

type UserMsgBroker interface {
	Created(ctx context.Context, event *model.Event) error
	Deleted(ctx context.Context, event *model.Event) error
	Updated(ctx context.Context, event *model.Event) error
}

type EventService interface {
	ProcessEvent
}
