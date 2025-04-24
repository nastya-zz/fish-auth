package broker

import (
	"auth/internal/model"
	"context"
)

type UserMsgBroker interface {
	Created(ctx context.Context, event model.Event) error
	Deleted(ctx context.Context, event model.Event) error
	Updated(ctx context.Context, event model.Event) error
}
