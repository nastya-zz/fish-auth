package auth

import (
	"auth/internal/utils"
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(utils.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	log.Printf("Claims %+v", claims)

	accessToken, err := i.authService.GetAccessToken(ctx, claims)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to get access token")
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
