package model

type Event struct {
	ID      int    `db:"id"`
	Type    string `db:"event_type"`
	Payload []byte `db:"payload"`
}

const (
	UserCreate = "user_create"
	UserUpdate = "user_update"
	UserDelete = "user_delete"
)
