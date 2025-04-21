package auth

import (
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	email := req.GetLogin()
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "login or password is empty")
	}

	token, err := i.authService.Login(ctx, email, password)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user not found")
	}

	return &desc.LoginResponse{
		RefreshToken: token,
	}, nil
}
