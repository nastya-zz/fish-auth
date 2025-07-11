package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/user-v1.UserV1/Get"
	ChatPath    = "/chat_v1.ChatV1/SendMessage"
)

const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type UserClaims struct {
	jwt.StandardClaims
	Name string `json:"name"`
	Role string `json:"role"`
	ID   string `json:"id"`
}
