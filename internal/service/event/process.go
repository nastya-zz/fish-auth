package event

import (
	"auth/internal/model"
	"context"
	"log/slog"
	"time"
)

const (
	eventBatchSize = 20
)

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

			events, err := s.eventRepository.GetNewEvent(ctx, eventBatchSize)
			if err != nil {
				log.Error("failed to get new event", err)
				continue
			}
			if events == nil {
				log.Debug("no new events")
				continue
			}

			s.SendMessage(ctx, events)
		}
	}()
}

func (s *Sender) SendMessage(ctx context.Context, events []*model.Event) {
	const op = "services.event-sender.SendMessage"

	for _, event := range events {
		log := slog.With(slog.String("op", op))
		log.Info("sending message", slog.Any("event", event))

		err := s.broker.Created(ctx, event)
		if err != nil {
			log.Error("failed to send message", err)
		}

		if err := s.eventRepository.SetDone(ctx, event.ID); err != nil {
			log.Error("failed to set event done", err)
		}
	}

}
