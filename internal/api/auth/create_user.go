package auth

import (
	"auth/internal/converter"
	errorsMsg "auth/pkg/api-errors-msg"
	"auth/pkg/logger"
	"context"
	"strings"

	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	logger.Info("Received user CreateRequest", "user_info", req.GetUserInfo())

	errors := validateCreateUserRequest(req)
	if len(errors) > 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"%s", strings.Join(errors, ", "))
	}
	preparedUser := converter.ToCreateUserFromDesc(req.GetUserInfo())

	id, err := i.authService.Create(ctx, preparedUser)

	logger.Info("Created user", "id", id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"%s", errorsMsg.UserCreationFailed)

	}

	return &desc.CreateUserResponse{Id: id}, nil
}

func validateCreateUserRequest(req *desc.CreateUserRequest) []string {
	user := req.GetUserInfo()
	var errors = make([]string, 0)

	if user == nil {
		errors = append(errors, errorsMsg.UserCreationFailed)
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
