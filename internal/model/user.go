package model

import (
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
	ID       string
	Password string
	Email    string
}

type UserPublish struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
}
