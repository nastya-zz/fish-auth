package utils

import (
	"auth/internal/model"
	errorsMsg "auth/pkg/api-errors-msg"
	"context"
	"google.golang.org/grpc/metadata"
	"strings"
)

const authorizationHeader = "authorization"

func GetClaims(ctx context.Context) (*model.UserClaims, string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errorsMsg.MetadataNotProvided
	}

	authHeader, ok := md[authorizationHeader]
	if !ok || len(authHeader) == 0 {
		return nil, errorsMsg.AuthorizationHeaderNotProvided
	}

	if !strings.HasPrefix(authHeader[0], AuthPrefix) {
		return nil, errorsMsg.AuthorizationHeaderInvalid
	}

	accessToken := strings.TrimPrefix(authHeader[0], AuthPrefix)

	claims, err := VerifyToken(accessToken, []byte(AccessTokenSecretKey))
	if err != nil {
		return nil, errorsMsg.AuthInvalidAccessToken
	}
	return claims, ""
}
