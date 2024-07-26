package models

import (
	"github.com/dgrijalva/jwt-go"
)

// User represents a User schema
type User struct {
	Base
	Email	string `json:"email" gorm:"unique;not null"`
	CurrentShopID *string `json:"current_shop_id" gorm:"null"`
	// CurrentShop Shop `json:"current_shop" gorm:"foreignKey:id;references:current_shop_id;save_associations:false"`
	CurrentShop   Shop   `json:"current_shop" gorm:"foreignKey:CurrentShopID;references:ID"`
}

// UserErrors represent the error format for user routes
type UserErrors struct {
	Err      bool   `json:"error"`
	Email    string `json:"email"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
}
