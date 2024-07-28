package models

import (
	pg "github.com/lib/pq"
)

// i think the table below is a waste of space
// and time. keeping it around for now. unnecessary complexity
type Notification struct {
	Base
	ShopID string `json:"shop_id"`
	NotificationType string `json:"notification_type"`
	Configured bool `json:"configured"`
}

// everytime a campaign is launched, a new record is created
type NotificationCampaign struct {
	Base
	// why not lol, it's not like i have any users
	ShopID string `json:"shop_id"`
	Shop *Shop `json:"shop" gorm:"foreignKey:ShopID;references:ID;"`
	
	NotificationConfigurationID string `json:"notification_configuration_id"`
	NotificationConfiguration *NotificationConfiguration `json:"notification_configuration" gorm:"foreignKey:NotificationConfigurationID;references:ID;"`
}

type NotificationsSent struct {
	Base

	// quick access
	NotificationCampaignID string `json:"notification_campaign_id"`
	NotificationCampaign *NotificationCampaign `json:"notification_campaign" gorm:"foreignKey:NotificationCampaignID;references:ID;"`

	// for identifying issues
	Status string `json:"status"` // sent, failed, pending
	APIResponse string `json:"api_response"` // headers in the case of web push
	APIStatus int `json:"api_status"`
}


// primarily just for Web Push Notifications
type NotificationConfiguration struct {
	Base
	NotificationType string `json:"notification_type"` // email, push, sms
	ShopID string `json:"shop_id"`
	Title string `json:"title"`
	Message string `json:"message"`
	URL string `json:"url"`
	Icon string `json:"icon"`
	Badge string `json:"badge"`
}

type NotificationSubscription struct {
	Base
	// endpoint can be used as a unique identifier
	// to track the user: https://stackoverflow.com/a/63769192/12674948
	Endpoint string `json:"endpoint" gorm:"unique;null"`

	Auth string `json:"auth" gorm:"null"`
	P256dh string `json:"p256dh" gorm:"null"`

	// this will become the shopify shop id soon
	OwnerID string `json:"owner_id"` // only for test notifications, from owner

	CustomerIDs pg.StringArray `json:"customer_ids" gorm:"type:text[]"`
	
	ShopID string `json:"shop_id" gorm:"null"`
	// Shop Shop `json:"shop" gorm:"foreignKey:id;references:shop_id;"`
	Shop Shop `json:"shop" gorm:"foreignKey:ShopID;references:ID;null"`
}
