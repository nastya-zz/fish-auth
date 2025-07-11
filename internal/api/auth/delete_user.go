package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	errorsMsg "auth/pkg/api-errors-msg"
	"auth/pkg/logger"
	"context"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DeleteUser(ctx context.Context, req *desc.BlockUserRequest) (*desc.BlockUserResponse, error) {
	const op = "api.auth.DeleteUser"

	userId := req.GetId()
	if len(userId) == 0 {
		return nil, status.Error(codes.InvalidArgument, errorsMsg.UserNotFound)
	}

	claims, strErr := utils.GetClaims(ctx)
	if strErr != "" {
		return nil, status.Errorf(codes.Aborted, "%s", strErr)
	}

	if userId != claims.ID && claims.Role != model.RoleAdmin {
		logger.Info(op, "userId", userId, "claims.ID", claims.ID, "equals", userId == claims.ID)
		return nil, status.Error(codes.PermissionDenied, errorsMsg.PermissionsDenied)
	}

	err := i.authService.Delete(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, errorsMsg.UserDeleteFailed)
	}

	return &desc.BlockUserResponse{Id: userId}, nil
}
