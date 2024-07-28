package models

import (
	"github.com/lib/pq"
)

// retire this model
// pretend it doesn't exist
type TrackedUser struct {
	Base
	FingerPrint string `json:"fingerprint"`
	ShopifyUniqueIDs pq.StringArray `json:"shopify_unique_ids" gorm:"type:text[]"`
	SubscriptionMetadata pq.StringArray `json:"subscription_metadata" gorm:"type:text[]"`
	Stores pq.StringArray `json:"stores" gorm:"type:text[]"`
}

type TrackedClick struct {
	Base
	NotificationCampaignID string `json:"notification_campaign_id"`
	NotificationCampaign NotificationCampaign `json:"notification_campaign"`

	ClickHeaders string `json:"click_headers"`
	ClickIP string `json:"click_ip"`

	// just in case NotificationCampaign.NotificationConfiguration.URL is different
	ClickForwardURL string `json:"click_forward_url"`

	// for quick indexed access
	ClickReferrer string `json:"click_referrer"`
	ClickUserAgent string `json:"click_user_agent"`
	ClickOS string `json:"click_os"`
}
