package auth

import (
	"context"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth/internal/utils"
	errorsMsg "auth/pkg/api-errors-msg"
	"auth/pkg/logger"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	logger.Info("GetAccessToken", "refreshToken", req.GetRefreshToken())
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(utils.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, errorsMsg.AuthInvalidRefreshToken)
	}

	logger.Info("Claims", "claims", claims)

	accessToken, err := i.authService.GetAccessToken(ctx, claims)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errorsMsg.JwtTokenFailed)
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
