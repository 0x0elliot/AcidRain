package models

import (
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// User represents a User schema
type User struct {
	Base
	Email	string `json:"email" gorm:"unique;not null"`
	CurrentShopID string `json:"current_shop_id"`
	CurrentShop Shop `json:"current_shop" gorm:"foreignKey:CurrentShopID;references:ID"`
}

// UserErrors represent the error format for user routes
type UserErrors struct {
	Err      bool   `json:"error"`
	Email    string `json:"email"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uuid.UUID `gorm:"primaryKey"`
}
