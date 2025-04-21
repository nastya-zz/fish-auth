package auth

import (
	"auth/internal/client/db"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

func (s serv) Check(ctx context.Context, endpoint string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(
	authRepository repository.AuthRepository,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}
