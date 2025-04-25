package event

import (
	"auth/internal/model"
	"auth/internal/service"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Broker struct {
	ch *amqp.Channel
}

func NewBroker(channel *amqp.Channel) service.UserMsgBroker {
	return &Broker{
		ch: channel,
	}
}

func (s Broker) Created(ctx context.Context, event *model.Event) error {
	return s.publish(ctx, "User.Created", event)
}

func (s *Broker) Deleted(ctx context.Context, event *model.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s *Broker) Updated(ctx context.Context, event *model.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s Broker) publish(_ context.Context, routingKey string, event *model.Event) error {
	const op = "broker.publish"
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		return fmt.Errorf("could not encode event: %w", err)
	}
	err := s.ch.Publish(
		"user_events", // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			AppId:       "auth_grpc_server",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		return fmt.Errorf("could not publish: %s, %w", op, err)
	}

	return nil
}
