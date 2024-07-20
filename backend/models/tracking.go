package models

import (
	"github.com/lib/pq"
)


type TrackedUser struct {
	Base
	FingerPrint string `json:"fingerprint"`
	ShopifyUniqueIDs pq.StringArray `json:"shopify_unique_ids" gorm:"type:text[]"`
	SubscriptionMetadata pq.StringArray `json:"subscription_metadata" gorm:"type:text[]"`
	Stores pq.StringArray `json:"stores" gorm:"type:text[]"`
}
