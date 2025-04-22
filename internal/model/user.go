package model

import (
	"time"
)

type User struct {
	ID         string
	Name       string
	Email      string
	Role       string
	Password   string
	IsVerified bool
	IsBlocked  bool
	CreatedAt  time.Time
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
