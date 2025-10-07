package auth

import (
	"context"
	"strings"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth/internal/converter"
	errorsMsg "auth/pkg/api-errors-msg"
)

func (i *Implementation) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.UpdateUserResponse, error) {
	errors := validateUpdateUserRequest(req)
	if len(errors) > 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			strings.Join(errors, ", "))
	}
	preparedUser := converter.ToUpdateUserFromDesc(req.GetId(), req.GetUserInfo())
	updatedUser, err := i.authService.UpdateUser(ctx, preparedUser)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			errorsMsg.UserUpdateFailed)
	}

	return &desc.UpdateUserResponse{
		Id:         updatedUser.ID,
		Name:       updatedUser.Name,
		Email:      updatedUser.Email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

func validateUpdateUserRequest(req *desc.UpdateUserRequest) []string {
	user := req.GetUserInfo()
	var errors = make([]string, 0)

	if user == nil {
		errors = append(errors, errorsMsg.UserUpdateFailed)
		return errors
	}
	if user.GetName() == "" {
		errors = append(errors, errorsMsg.UsernameInvalid)
	}
	if user.GetEmail() == "" {
		errors = append(errors, errorsMsg.UserEmailInvalid)
	}
	if user.GetPassword() == "" {
		errors = append(errors, errorsMsg.UserPasswordInvalid)
	}

	password := user.GetPassword()
	confirmPass := user.GetPasswordConfirm()

	if password != confirmPass {
		errors = append(errors, errorsMsg.UserPasswordConfirmInvalid)
	}

	return errors
}
