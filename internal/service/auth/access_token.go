package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"fmt"
)

func (s serv) GetAccessToken(ctx context.Context, claims *model.UserClaims) (string, error) {
	const op = "auth/service.GetAccessToken"

	accessToken, err := utils.GenerateToken(model.User{
		Name: claims.Name,
		Role: claims.Role,
		ID:   claims.ID,
	}, []byte(utils.AccessTokenSecretKey), utils.AccessTokenExpiration)

	if err != nil {
		return "", fmt.Errorf(op+": %w", err)
	}

	return accessToken, nil
}
