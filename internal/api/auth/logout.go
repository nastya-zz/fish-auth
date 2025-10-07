package auth

import (
	"context"
	"strings"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"auth/internal/utils"
	errorsMsg "auth/pkg/api-errors-msg"
	"auth/pkg/logger"
)

//todo: добавить инвалидацию refresh токена

func (i *Implementation) Logout(ctx context.Context, req *desc.LogoutRequest) (*desc.LogoutResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, errorsMsg.MetadataNotProvided)
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, errorsMsg.AuthorizationHeaderNotProvided)
	}

	if !strings.HasPrefix(authHeader[0], utils.AuthPrefix) {
		return nil, status.Errorf(codes.InvalidArgument, errorsMsg.AuthorizationHeaderInvalid)
	}

	tokenStr := strings.TrimPrefix(authHeader[0], utils.AuthPrefix)
	logger.Info("Logout", "token", tokenStr)

	claims, err := utils.VerifyToken(tokenStr, []byte(utils.AccessTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, errorsMsg.AuthInvalidAccessToken)
	}

	utils.RevokeToken(tokenStr, claims.ExpiresAt)
	return &desc.LogoutResponse{
		Success: true,
	}, nil
}
