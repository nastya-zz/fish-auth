package auth

import (
	"auth/internal/repository/auth/converter"
	apierrors "auth/pkg/api-errors-msg"
	service_errors "auth/pkg/service-errors"
	"context"
	"errors"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	email := req.GetLogin()
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, apierrors.AuthInvalidParams)
	}

	token, role, err := i.authService.Login(ctx, email, password)
	if err != nil {
		parsedErr := getError(err)
		return nil, status.Errorf(codes.Aborted, parsedErr)
	}

	return &desc.LoginResponse{
		RefreshToken: token,
		Role:         converter.DescRole(role),
	}, nil
}

func getError(err error) string {
	switch true {
	case errors.Is(err, service_errors.PasswordNotMatched):
	case errors.Is(err, service_errors.RefreshTokenFailed):
	case errors.Is(err, service_errors.UserBlocked):
		return err.Error()
	}
	return "Произошла ошибка при попытке авторизоваться"
}
