package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
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
	ID    int64
}
