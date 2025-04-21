package auth

import (
	"auth/internal/utils"
	"context"
	"errors"
	"fmt"
)

func (s serv) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth/auth.Login"

	user, err := s.authRepository.Login(ctx, email)
	if err != nil {
		return "", fmt.Errorf(op+": %s, %s, %w", email, err)
	}

	equal := utils.VerifyPassword(user.Password, password)
	if !equal {
		return "", errors.New("passwords do not match")
	}

	refreshToken, err := utils.GenerateToken(*user, []byte(utils.RefreshTokenSecretKey), utils.RefreshTokenExpiration)
	if err != nil {
		return "", fmt.Errorf(op+": refreshToken: %s, %w", email, err)
	}

	return refreshToken, nil
}
