package event

import (
	"auth/internal/repository"
	"auth/internal/service"
)

type Sender struct {
	eventRepository repository.EventRepository
	broker          service.UserMsgBroker
}

func NewService(storage repository.EventRepository, broker service.UserMsgBroker) service.EventService {
	return &Sender{
		eventRepository: storage,
		broker:          broker,
	}
}
