package models

type Shop struct {
	Base
	Name string `json:"name" gorm:"not null"`
	ShopIdentifier string `json:"shop_identifier" gorm:"unique;not null"`
	AccessToken string `json:"access_token"`
	Platform string `json:"platform"` // shopify
	OwnerID string `json:"owner_id"`

	// Replace struct Notification by entries here in the near future
}
