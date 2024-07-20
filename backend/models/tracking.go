package models

type TrackedUser struct {
	Base
	FingerPrint string `json:"fingerprint"`
	ShopifyUniqueIDs []string `json:"shopify_unique_ids" gorm:"type:text[]"`
	SubscriptionMetadata []string `json:"subscription_metadata" gorm:"type:text[]"`
	Stores []string `json:"stores" gorm:"type:text[]"`
}
