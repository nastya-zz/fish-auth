package event

import (
	"auth/internal/model"
	"auth/pkg/logger"
	"context"
	"time"
)

const (
	eventBatchSize = 20
)

func (s *Sender) StartProcessEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "services.event-sender.StartProcessEvents"

	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("stopping event processing", "op", op)
				return
			case <-ticker.C:
				// noop
			}

			events, err := s.eventRepository.GetNewEvent(ctx, eventBatchSize)
			if err != nil {
				logger.Error("failed to get new event", "error", err, "op", op)
				continue
			}
			if events == nil {
				logger.Debug("no new events", "op", op)
				continue
			}

			s.SendMessage(ctx, events)
		}
	}()
}

func (s *Sender) SendMessage(ctx context.Context, events []*model.Event) {
	const op = "services.event-sender.SendMessage"

	for _, event := range events {
		logger.Info("sending message", "event", event, "op", op)

		err := s.broker.Created(ctx, event)
		if err != nil {
			logger.Error("failed to send message", "error", err, "op", op)
		}

		if err := s.eventRepository.SetDone(ctx, event.ID); err != nil {
			logger.Error("failed to set event done", "error", err, "op", op)
		}
	}

}
