package event

import (
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
	"log/slog"
	"time"
)

type Sender struct {
	eventRepository repository.EventRepository
}

func NewService(storage repository.EventRepository) service.EventService {
	return &Sender{
		eventRepository: storage,
	}
}

func (s *Sender) StartProcessEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "services.event-sender.StartProcessEvents"

	log := slog.With(slog.String("op", op))

	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("stopping event processing")
				return
			case <-ticker.C:
				// noop
			}

			event, err := s.eventRepository.GetNewEvent(ctx)
			if err != nil {
				log.Error("failed to get new event", err)
				continue
			}
			if event == nil {
				log.Debug("no new events")
				continue
			}

			s.SendMessage(ctx, event)

			if err := s.eventRepository.SetDone(ctx, event.ID); err != nil {
				log.Error("failed to set event done", err)
			}
		}
	}()
}

func (s *Sender) SendMessage(ctx context.Context, event *model.Event) {
	const op = "services.event-sender.SendMessage"

	log := slog.With(slog.String("op", op))
	log.Info("sending message", slog.Any("event", event))

	// TODO: implement sending message to the external service.

}
