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
