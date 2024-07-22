package models

import (
	pg "github.com/go-pg/pg/v10"
)

// i think the table below is a waste of space
// and time. keeping it around for now. unnecessary complexity
type Notification struct {
	Base
	ShopID string `json:"shop_id"`
	NotificationType string `json:"notification_type"`
	Configured bool `json:"configured"`
}

type NotificationsSent struct {
	Base
	NotificationID string `json:"notification_id"`
	NotificationSubscriptionID string `json:"notification_subscription_id"`
	NotificationStatus string `json:"notification_status"`
}

type NotificationSubscription struct {
	Base
	// endpoint can be used as a unique identifier
	// to track the user: https://stackoverflow.com/a/63769192/12674948
	Endpoint string `json:"endpoint" gorm:"unique;null"`

	Auth string `json:"auth" gorm:"null"`
	P256dh string `json:"p256dh" gorm:"null"`

	// this will become the shopify shop id soon
	OwnerID string `json:"owner_id"` // only for test notifications

	CustomerIDs pg.IntSlice `json:"customer_ids" gorm:"type:int[]"`
	
	// might retire these soon
	NotificationID string `json:"notification_id"`
	ShopifyUniqueID string `json:"shopify_unique_id"`
}
