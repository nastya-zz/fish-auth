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
func (UnimplementedAuthV1Server) UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedAuthV1Server) Check(context.Context, *CheckRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}

func (UnimplementedAuthV1Server) DeleteUser(context.Context, *BlockUserRequest) (*BlockUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedAuthV1Server) UpdateUserRole(context.Context, *UpdateUserRoleRequest) (*UpdateUserRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserRole not implemented")
}
*/
