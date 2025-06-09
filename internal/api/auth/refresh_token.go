package auth

import (
	"auth/internal/utils"
	apierrors "auth/pkg/api-errors-msg"
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(utils.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, apierrors.AuthInvalidRefreshToken)
	}

	log.Printf("Claims %+v", claims)

	user, err := i.authService.GetUser(ctx, claims.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, apierrors.UserNotFound)
	}

	refreshToken, err := utils.GenerateToken(*user,
		[]byte(utils.RefreshTokenSecretKey),
		utils.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
