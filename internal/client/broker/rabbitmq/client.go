package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueName = "user"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewRabbitMQ(ctx context.Context, url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error connection %s: %w", "amqp.Dial", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error connection %s: %w", "conn.Channel", err)
	}

	_, err = channel.QueueDeclare(QueueName, true, false, false, false, nil)

	if err != nil {
		return nil, fmt.Errorf("error create queue  %s: %w", "ch.QueueDeclare", err)
	}

	if err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, fmt.Errorf("error connection  %s: %w", "ch.Qos", err)
	}

	// XXX: Dead Letter Exchange will be implemented in future episodes

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}, nil
}

// Close ...
func (r *RabbitMQ) Close() error {
	if err := r.Connection.Close(); err != nil {
		return fmt.Errorf("error connection  %s: %w", "Connection.Close", err)
	}

	return nil
}

func (r *RabbitMQ) Connect() *RabbitMQ {
	return r
}
