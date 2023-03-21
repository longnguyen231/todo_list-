package models

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type User struct {
	Id       int    `json:"id" :"id"`
	Name     string `json:"name" :"name"`
	Email    string `json:"email" :"email"`
	Password string `json:"password" :"password"`
	Age      int    `json:"age" :"age"`
	Avatar   string `json:"avatar" :"avatar"`
}
type JwtCustomClaims struct {
	Id                   int    `json:"id" :"id"`
	Email                string `json:"email" :"email"`
	Password             bool   `json:"password" :"password"`
	jwt.RegisteredClaims `:"jwt.RegisteredClaims"`
}
type Task struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Description string    `json:"description"`
	Completed   bool      `json:"complete"`
	CreateAt    time.Time `json:"create_at"`
}
