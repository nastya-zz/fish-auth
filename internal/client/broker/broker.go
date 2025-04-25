package broker

import "auth/internal/client/broker/rabbitmq"

type ClientMsgBroker interface {
	Connect() *rabbitmq.RabbitMQ
	Close() error
}
