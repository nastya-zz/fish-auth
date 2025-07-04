package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Password   string
	IsVerified bool `json:"isVerified"`
	IsBlocked  bool
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time
	LastLogin  time.Time
}

type CreateUser struct {
	Name     string
	Email    string
	Role     string
	Password string
}

type UpdateUser struct {
	ID         string `db:"id"`
	Name       string `db:"name"`
	Password   string
	Email      string `db:"email"`
	IsVerified bool   `db:"isVerified"`
}

type UserPublish struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
}

func GetUuid[T ~string](id T) (uuid.UUID, error) {
	return uuid.Parse(string(id))
}
