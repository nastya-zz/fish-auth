package auth

import (
	"auth/internal/converter"
	errorsMsg "auth/pkg/api-errors-msg"
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Printf("Received user CreateRequest %+v", req.GetUserInfo())

	errors := validateCreateUserRequest(req)
	if len(errors) > 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			strings.Join(errors, ", "))
	}
	preparedUser := converter.ToCreateUserFromDesc(req.GetUserInfo())

	id, err := i.authService.Create(ctx, preparedUser)

	log.Printf("Created user with id %s", id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			errorsMsg.UserCreationFailed)

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
