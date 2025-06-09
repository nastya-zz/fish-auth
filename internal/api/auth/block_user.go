package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	errorsMsg "auth/pkg/api-errors-msg"
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) BlockUser(ctx context.Context, req *desc.BlockUserRequest) (*desc.BlockUserResponse, error) {
	const op = "api.auth.BlockUser"

	userId := req.GetId()
	if len(userId) == 0 {
		return nil, status.Error(codes.InvalidArgument, errorsMsg.UserNotFound)
	}

	claims, strErr := utils.GetClaims(ctx)
	if strErr != "" {
		return nil, status.Errorf(codes.Aborted, strErr)
	}

	if claims.Role != model.RoleAdmin {
		return nil, status.Error(codes.PermissionDenied, errorsMsg.PermissionsDenied)
	}

	id, err := i.authService.BlockUser(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, errorsMsg.UserBlockFailed)

	}

	return &desc.BlockUserResponse{Id: id}, nil
}
