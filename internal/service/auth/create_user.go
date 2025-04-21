package auth

import (
	"auth/internal/model"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func (s serv) Create(ctx context.Context, user *model.CreateUser) (string, error) {

	password := user.Password
	hash, err := HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hash

	id, err := s.authRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
