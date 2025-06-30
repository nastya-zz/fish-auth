package main

import (
	"auth/internal/app"
	"auth/internal/logger"
	"context"
)

func main() {
	// Инициализируем логгер в самом начале
	logger.Init()
	
	ctx := context.Background()
	logger.Info("Starting auth service")

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to init app", "error", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		logger.Fatal("failed to run app", "error", err.Error())
	}
}
