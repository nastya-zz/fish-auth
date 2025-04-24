package auth

import (
	"auth/internal/client/db"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
)

type serv struct {
	authRepository  repository.AuthRepository
	eventRepository repository.EventRepository
	eventService    service.EventService
	txManager       db.TxManager
}

func (s serv) Check(ctx context.Context, endpoint string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(
	authRepository repository.AuthRepository,
	eventRepository repository.EventRepository,
	eventService service.EventService,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository:  authRepository,
		eventRepository: eventRepository,
		eventService:    eventService,
		txManager:       txManager,
	}
}
