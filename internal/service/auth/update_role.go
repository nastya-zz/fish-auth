package auth

import "context"

func (s serv) UpdateRole(ctx context.Context, id, role string) error {
	return s.authRepository.UpdateRole(ctx, id, role)
}
