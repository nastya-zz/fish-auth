package auth

import (
	"auth/internal/repository/auth/converter"
	"auth/pkg/logger"
	"context"

	"auth/internal/model"
	"auth/internal/utils"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateUserRole(ctx context.Context, req *desc.UpdateUserRoleRequest) (*desc.UpdateUserRoleResponse, error) {
	const op = "auth.api.UpdateUserRole"
	userID := req.GetId()
	claims, strErr := utils.GetClaims(ctx)

	if strErr != "" {
		return nil, status.Errorf(codes.Aborted, strErr)
	}

	if claims.Role != model.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}

	if userID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	role := req.Role.String()
	if role == "" {
		return nil, status.Errorf(codes.InvalidArgument, "role is required")
	}

	err := i.authService.UpdateRole(ctx, userID, role)
	if err != nil {
		logger.Error(op+"Error updating user role: ", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to update user role: %v", err)
	}

	return &desc.UpdateUserRoleResponse{
		Id:   userID,
		Role: converter.DescRole(role),
	}, nil
}
