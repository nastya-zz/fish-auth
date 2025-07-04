package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"auth/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
)

func (s serv) UpdateUser(ctx context.Context, user *model.UpdateUser) (*model.UpdateUser, error) {
	logger.Info("start to update user", "user_id", user.ID, "email")

	var hash string
	var err error
	var updatedUser *model.UpdateUser
	if user.Password != "" {
		password := user.Password
		hash, err = utils.HashPassword(password)
	}
	if err != nil {
		logger.Error("failed to hash password", "user_id", user.ID, "email", user.Email, "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hash

	needPublish := false

	existsUser, err := s.GetUser(ctx, user.ID)
	if err != nil {
		logger.Error("get user for update is failed", "user_id", user.ID, "error", err)
		return nil, fmt.Errorf("get user for update is failed: %w", err)
	}

	needPublish = user.IsVerified != existsUser.IsVerified || user.Email != existsUser.Email || user.Name != existsUser.Name

	if needPublish {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			var errTx error
			updatedUser, errTx = s.authRepository.Update(ctx, &model.UpdateUser{
				ID:         user.ID,
				Name:       user.Name,
				Email:      user.Email,
				Password:   user.Password,
				IsVerified: user.IsVerified,
			})
			if errTx != nil {
				return errTx
			}

			body, errTx := json.Marshal(model.UserPublish{
				ID:         user.ID,
				Name:       user.Name,
				Email:      user.Email,
				IsVerified: user.IsVerified,
			})

			if errTx != nil {
				return fmt.Errorf("error in marshal json body %w", errTx)
			}

			event := &model.Event{
				Type:    model.UserUpdate,
				Payload: body,
			}

			errTx = s.eventRepository.SaveEvent(ctx, event)
			if errTx != nil {
				return errTx
			}

			return nil
		})

		if err != nil {
			logger.Error("transaction failed during user update", "error", err, "user_id", user.ID, "email", user.Email)
			return nil, err
		}
	} else {
		updatedUser, err = s.authRepository.Update(ctx, &model.UpdateUser{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Password:   user.Password,
			IsVerified: user.IsVerified,
		})

		if err != nil {
			logger.Error("repository failed during user update", "error", err, "user_id", user.ID, "email", user.Email)
			return nil, err
		}
	}

	logger.Info("user updated successfully", "user_id", user.ID, "email", user.Email)

	return updatedUser, nil
}
