package service

import (
	"auth/internal/model"
	"context"
	"time"
)

type AuthService interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UpdateUser) (*model.UpdateUser, error)
	Login(ctx context.Context, login string, password string) (string, string, error)
	GetAccessToken(ctx context.Context, claims *model.UserClaims) (string, error)
	GetRefreshToken(ctx context.Context, claims *model.UserClaims) (string, error)
	Create(ctx context.Context, user *model.CreateUser) (string, error)
	Check(ctx context.Context, endpoint string) (bool, error)
	BlockUser(ctx context.Context, id string) (string, error)
	Delete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id, role string) error
}

type ProcessEvent interface {
	StartProcessEvents(ctx context.Context, handlePeriod time.Duration)
	SendMessage(ctx context.Context, event []*model.Event)
}

type UserMsgBroker interface {
	SendEvent(ctx context.Context, event *model.Event) error
}

type EventService interface {
	ProcessEvent
}
