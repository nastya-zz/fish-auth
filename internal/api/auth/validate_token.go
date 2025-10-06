package auth

import (
	"context"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth/internal/utils"
	errorsMsg "auth/pkg/api-errors-msg"
)

func (i *Implementation) ValidateToken(ctx context.Context, req *desc.ValidateTokenRequest) (*desc.ValidateTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetToken(), []byte(utils.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, errorsMsg.AuthInvalidRefreshToken)
	}

	return &desc.ValidateTokenResponse{Claims: &desc.Claims{
		Id:   claims.ID,
		Role: claims.Role,
	}}, nil
}
