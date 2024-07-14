package models

type Notification struct {
	Base
	ShopID string `json:"shop_id"`
	NotificationType string `json:"notification_type"`
	NotificationData string `json:"notification_data"`
	NotificationStatus string `json:"notification_status"`
}

type NotificationSubscription struct {
	Base
	Endpoint string `json:"endpoint"`

	Auth string `json:"auth"`
	P256dh string `json:"p256dh"`

	OwnerID string `json:"owner_id"`
}
