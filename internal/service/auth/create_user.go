package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"fmt"
)

func (s serv) Create(ctx context.Context, user *model.CreateUser) (string, error) {

	password := user.Password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hash

	id, err := s.authRepository.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}
