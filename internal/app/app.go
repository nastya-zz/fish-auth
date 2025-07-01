package app

import (
	"auth/internal/closer"
	"auth/internal/config"
	"auth/pkg/logger"
	"context"
	descAuth "github.com/nastya-zz/fisher-protocols/gen/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

const (
	envDev  = "dev"
	envProd = "prod"
)

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	a.runEventSender(ctx)

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info("GRPC server is running", "address", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		logger.Error("failed to listen", "error", err, "address", a.serviceProvider.GRPCConfig().Address())
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		logger.Error("failed to serve grpc", "error", err)
		return err
	}

	return nil
}

func (a *App) runEventSender(ctx context.Context) {
	a.serviceProvider.eventService.StartProcessEvents(ctx, 10*time.Second)
}
