package auth

import (
	"auth/internal/model"
	"context"
	"encoding/json"
	"fmt"
)

func (s serv) Delete(ctx context.Context, id string) error {
	const op = "service.auth.DeleteUser"

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.authRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		body, errTx := json.Marshal(model.UserPublish{
			ID: id,
		})

		if errTx != nil {
			return fmt.Errorf("error in marshal json body %w", errTx)
		}

		event := &model.Event{
			Type:    model.UserDelete,
			Payload: body,
		}

		errTx = s.eventRepository.SaveEvent(ctx, event)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		fmt.Printf("%s: %s", op, err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
