package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"encoding/json"
	"fmt"
)

func (s serv) Create(ctx context.Context, user *model.CreateUser) (string, error) {

	password := user.Password
	hash, err := utils.HashPassword(password)
	if err != nil {
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
		return "", err
	}

	return id, nil
}
