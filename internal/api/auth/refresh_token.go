package auth

import (
	"auth/internal/utils"
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(utils.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	log.Printf("Claims %+v", claims)
	return nil, status.Errorf(codes.Unimplemented, "method GetRefreshToken not implemented")
}
