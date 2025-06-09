package auth

import (
	"auth/internal/model"
	"context"
)

func (s serv) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.authRepository.Get(ctx, id)
}
