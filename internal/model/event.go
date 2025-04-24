package model

type Event struct {
	ID      int
	Type    string
	Payload []byte
}

const (
	UserCreate = "user_create"
	UserUpdate = "user_update"
	UserDelete = "user_delete"
)
