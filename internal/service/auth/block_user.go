package auth

import (
	"context"
	"fmt"
)

func (s serv) BlockUser(ctx context.Context, id string) (string, error) {
	const op = "auth.BlockUser"

	if err := s.authRepository.Block(ctx, id); err != nil {
		return "", fmt.Errorf(op+"%w", err)
	}
	return id, nil
}
