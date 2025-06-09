package auth

import (
	"auth/internal/utils"
	service_errors "auth/pkg/service-errors"
	"context"
	"fmt"
)

func (s serv) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth/auth.Login"

	user, err := s.authRepository.Login(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: %s, %w", op, email, err)
	}

	if user.IsBlocked {
		return "", service_errors.UserBlocked
	}

	equal := utils.CheckPasswordHash(user.Password, password)
	if !equal {
		return "", service_errors.PasswordNotMatched
	}

	refreshToken, err := utils.GenerateToken(*user, []byte(utils.RefreshTokenSecretKey), utils.RefreshTokenExpiration)
	if err != nil {
		return "", service_errors.RefreshTokenFailed
	}

	return refreshToken, nil
}
