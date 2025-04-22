package model

import (
	"database/sql"
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
	LastLogin  sql.NullTime
}

type CreateUser struct {
	Name     string
	Email    string
	Role     string
	Password string
}

type UpdateUser struct {
	Name  string
	Email string
	ID    string
}
