package auth

import (
	"auth/internal/utils"
	"auth/pkg/logger"
	errorsMsg "auth/pkg/api-errors-msg"
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
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
