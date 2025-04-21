package auth

import (
	"auth/internal/client/db"
	"auth/internal/repository"
	"auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
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
