package models

import (
	"github.com/dgrijalva/jwt-go"
)

// User represents a User schema
type User struct {
	Base
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Username string `json:"username"`
}

// UserErrors represent the error format for user routes
type UserErrors struct {
	Err      bool   `json:"error"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey"`
}
