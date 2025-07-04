package model

import (
	"database/sql"
	"github.com/google/uuid"
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
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	IsVerified bool      `db:"isVerified"`
}
