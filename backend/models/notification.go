package models

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
	Endpoint string `json:"endpoint"`

	Auth string `json:"auth"`
	P256dh string `json:"p256dh"`

	OwnerID string `json:"owner_id"` // only for test notifications
	NotificationID string `json:"notification_id"`
}
