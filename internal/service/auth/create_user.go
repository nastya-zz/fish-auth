package auth

import (
	"auth/pkg/logger"
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"encoding/json"
	"fmt"
)

func (s serv) Create(ctx context.Context, user *model.CreateUser) (string, error) {
	logger.Info("starting user creation", "email", user.Email, "role", user.Role)

	password := user.Password
	hash, err := utils.HashPassword(password)
	if err != nil {
		logger.Error("failed to hash password", "error", err, "email", user.Email)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hash

	var id string
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		createdUser, errTx := s.authRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}
		id = createdUser.ID

		body, errTx := json.Marshal(model.UserPublish{
			ID:         createdUser.ID,
			Name:       createdUser.Name,
			Email:      createdUser.Email,
			IsVerified: createdUser.IsVerified,
			CreatedAt:  createdUser.CreatedAt,
		})

		if errTx != nil {
			return fmt.Errorf("error in marshal json body %w", errTx)
		}

		event := &model.Event{
			Type:    model.UserCreate,
			Payload: body,
		}

		errTx = s.eventRepository.SaveEvent(ctx, event)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		logger.Error("transaction failed during user creation", "error", err, "email", user.Email)
		return "", err
	}

	logger.Info("user created successfully", "user_id", id, "email", user.Email)
	return id, nil
}
