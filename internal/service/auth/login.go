package auth

import (
	"auth/pkg/logger"
	"auth/internal/utils"
	service_errors "auth/pkg/service-errors"
	"context"
	"fmt"
)

func (s serv) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth/auth.Login"
	logger.Info("user login attempt", "email", email)

	user, err := s.authRepository.Login(ctx, email)
	if err != nil {
		logger.Error("failed to get user for login", "error", err, "email", email)
		return "", fmt.Errorf("%s: %s, %w", op, email, err)
	}

	if user.IsBlocked {
		logger.Warn("blocked user attempted login", "email", email, "user_id", user.ID)
		return "", service_errors.UserBlocked
	}

	equal := utils.CheckPasswordHash(user.Password, password)
	if !equal {
		logger.Warn("invalid password for user", "email", email)
		return "", service_errors.PasswordNotMatched
	}

	refreshToken, err := utils.GenerateToken(*user, []byte(utils.RefreshTokenSecretKey), utils.RefreshTokenExpiration)
	if err != nil {
		logger.Error("failed to generate refresh token", "error", err, "email", email)
		return "", service_errors.RefreshTokenFailed
	}

	logger.Info("user logged in successfully", "email", email, "user_id", user.ID)
	return refreshToken, nil
}
