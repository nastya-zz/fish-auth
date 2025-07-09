package auth

import (
	"auth/internal/service"
	desc "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}

/*
TODO:
 - реализовать проверку токенов для сторонних сервисов
 - реализовать провероку доступов до сервисов исходя из ролевой модели

func (UnimplementedAuthV1Server) Check(context.Context, *CheckRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
*/
