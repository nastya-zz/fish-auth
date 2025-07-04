package app

import (
	"auth/internal/api/auth"
	"auth/internal/client/broker"
	"auth/internal/client/broker/rabbitmq"
	"auth/internal/client/db"
	"auth/internal/client/db/pg"
	"auth/internal/closer"
	"auth/internal/config"
	"auth/internal/repository"
	authRepository "auth/internal/repository/auth"
	eventRepository "auth/internal/repository/event"
	"auth/internal/service"
	authService "auth/internal/service/auth"
	eventService "auth/internal/service/event"
	"auth/internal/transaction"
	"auth/pkg/logger"
	"context"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	rmqConfig  config.RMQConfig

	rmqClient       broker.ClientMsgBroker
	dbClient        db.Client
	txManager       db.TxManager
	authRepository  repository.AuthRepository
	eventRepository repository.EventRepository

	authService  service.AuthService
	eventService service.EventService
	eventBroker  service.UserMsgBroker

	authImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatal("failed to get pg config", "error", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) RMQConfig() config.RMQConfig {
	if s.rmqConfig == nil {
		cfg, err := config.NewRMQConfig()
		if err != nil {
			logger.Fatal("failed to get rmqConfig", "error", err.Error())
		}

		s.rmqConfig = cfg
	}

	return s.rmqConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			logger.Fatal("failed to get grpc config", "error", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create db client", "error", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("ping error", "error", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) RabbitMQClient(ctx context.Context) broker.ClientMsgBroker {
	if s.rmqClient == nil {
		cl, err := rabbitmq.NewRabbitMQ(ctx, s.RMQConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create rmq client", "error", err)
		}

		closer.Add(cl.Close)

		s.rmqClient = cl
	}

	return s.rmqClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) EventRepository(ctx context.Context) repository.EventRepository {
	if s.eventRepository == nil {
		s.eventRepository = eventRepository.NewRepository(s.DBClient(ctx))
	}

	return s.eventRepository
}

func (s *serviceProvider) EventBroker(ctx context.Context) service.UserMsgBroker {
	if s.eventBroker == nil {
		con := s.RabbitMQClient(ctx).Connect()
		s.eventBroker = eventService.NewBroker(con.Channel)
	}

	return s.eventBroker
}

func (s *serviceProvider) EventService(ctx context.Context) service.EventService {
	if s.eventService == nil {
		s.eventService = eventService.NewService(
			s.EventRepository(ctx),
			s.EventBroker(ctx),
		)
	}

	return s.eventService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthRepository(ctx),
			s.EventRepository(ctx),
			s.EventService(ctx),
			s.TxManager(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {

	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
