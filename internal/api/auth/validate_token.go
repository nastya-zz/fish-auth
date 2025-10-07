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

func (i *Implementation) ValidateToken(ctx context.Context, req *desc.ValidateTokenRequest) (*desc.ValidateTokenResponse, error) {
	logger.Info("ValidateToken", "token", req.GetToken())

	claims, err := utils.VerifyToken(req.GetToken(), []byte(utils.AccessTokenSecretKey))
	if err != nil {
		logger.Error("Error verifying token", "error", err)
		return nil, status.Errorf(codes.Aborted, errorsMsg.AuthInvalidRefreshToken)
	}

	logger.Info("Claims", "claims", claims)

	return &desc.ValidateTokenResponse{Claims: &desc.Claims{
		Id:   claims.ID,
		Role: claims.Role,
	}}, nil
}
