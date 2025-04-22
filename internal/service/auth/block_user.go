package auth

import "context"

func (s serv) BlockUser(ctx context.Context, id string) (string, error) {
	const op = "auth.BlockUser"

	if err := s.BlockUser(ctx, id); err != nil {
		return "", err
	}
	return id, nil
}
