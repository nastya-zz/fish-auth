package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"fmt"
)

func (s serv) GetRefreshToken(ctx context.Context, claims *model.UserClaims) (string, error) {
	const op = "auth/service.GetRefreshToken"

	refreshToken, err := utils.GenerateToken(model.User{
		Name: claims.Name,
		Role: claims.Role,
		ID:   claims.ID,
	}, []byte(utils.RefreshTokenSecretKey), utils.RefreshTokenExpiration)

	if err != nil {
		return "", fmt.Errorf(op+": %w", err)
	}

	return refreshToken, nil
}
