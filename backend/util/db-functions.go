package util

import (
	db "go-authentication-boilerplate/database"
	models "go-authentication-boilerplate/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUserById(id string) (*models.User, error) {
	user := new(models.User)
	txn := db.DB.Where("id = ?", id).First(&user)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting user: %v", txn.Error)
		return nil, txn.Error
	}
	return user, nil
}

func GetNotificationSubscriptionById(id string) (*models.NotificationSubscription, error) {
	subscription := new(models.NotificationSubscription)
	txn := db.DB.Where("id = ?", id).First(&subscription)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting subscription: %v", txn.Error)
		return nil, txn.Error
	}
	return subscription, nil
}

func SetNotificationConfiguration(notificationConfiguration *models.NotificationConfiguration) (*models.NotificationConfiguration, error) {
	// check if notificationConfiguration with ID exists
	if notificationConfiguration.ID == "" {
		notificationConfiguration.ID = uuid.New().String()
		notificationConfiguration.CreatedAt = db.DB.NowFunc().String()
		notificationConfiguration.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Create(notificationConfiguration)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating notification configuration: %v", txn.Error)
			return notificationConfiguration, txn.Error
		}
	} else {
		notificationConfiguration.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Save(notificationConfiguration)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving notification configuration: %v", txn.Error)
			return notificationConfiguration, txn.Error
		}
	}

	return notificationConfiguration, nil
}

func GetCountOfNotificationSubscriptionsByShopId(shopId string) (int64, error) {
	var count int64
	// Don't count the notifications with a owner_id
	txn := db.DB.Model(&models.NotificationSubscription{}).Where("shop_id = ? AND owner_id = ''", shopId).Count(&count)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting subscription count: %v", txn.Error)
		return count, txn.Error
	}
	return count, nil
}

func GetNotificationSubscriptionByShopId(shopId string) ([]models.NotificationSubscription, error) {
	subscriptions := []models.NotificationSubscription{}
	txn := db.DB.Where("shop_id = ?", shopId).Find(&subscriptions)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting subscription: %v", txn.Error)
		return subscriptions, txn.Error
	}
	return subscriptions, nil
}

// exclusively for test notifications
func GetNoficationSubscriptionByOwnerId(ownerID string) ([]models.NotificationSubscription, error) {
	subscription := models.NotificationSubscription{}
	txn := db.DB.Where("owner_id = ?", ownerID).Find(&subscription)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting subscription: %v", txn.Error)
		return nil, txn.Error
	}

	return []models.NotificationSubscription{subscription}, nil
}

func GetShops(ownerID string) ([]models.Shop, error) {
	shops := []models.Shop{}
	txn := db.DB.Where("owner_id = ?", ownerID).Find(&shops)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting shops: %v", txn.Error)
		return shops, txn.Error
	}
	return shops, nil
}

func GetIdenticalSubscription(subscription models.NotificationSubscription) (*models.NotificationSubscription, error) {
	sub := new(models.NotificationSubscription)
	txn := db.DB.Where("endpoint = ? AND auth = ? AND p256dh = ?", subscription.Endpoint, subscription.Auth, subscription.P256dh).First(sub)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting subscription: %v", txn.Error)
		return sub, txn.Error
	}
	return sub, nil
}

func GetSubscriptionFromEndpoint(endpoint string) (*models.NotificationSubscription, error) {
	sub := new(models.NotificationSubscription)
	txn := db.DB.Where("endpoint = ?", endpoint).First(sub)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting subscription: %v", txn.Error)
		return sub, txn.Error
	}
	return sub, nil
}

func GetShopById(shopId string) (*models.Shop, error) {
	shop := new(models.Shop)
	txn := db.DB.Where("id = ?", shopId).First(shop)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting shop: %v", txn.Error)
		return shop, txn.Error
	}
	return shop, nil
}

func GetNotificationConfigurationById(id string, shopId string) (*models.NotificationConfiguration, error) {
	notificationConfiguration := new(models.NotificationConfiguration)
	txn := db.DB.Where("id = ? AND shop_id = ?", id, shopId).First(notificationConfiguration)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting notification configuration: %v", txn.Error)
		return notificationConfiguration, txn.Error
	}
	return notificationConfiguration, nil
}

func GetNotificationConfigurationsById(id string) ([]*models.NotificationConfiguration, error) {
	notificationConfigurations := []*models.NotificationConfiguration{}
	txn := db.DB.Where("shop_id = ?", id).Find(&notificationConfigurations)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting campaigns: %v", txn.Error)
		return notificationConfigurations, txn.Error
	}
	return notificationConfigurations, nil
}	

func GetNotificationSentById(id string) (*models.NotificationsSent, error) {
	notificationSent := new(models.NotificationsSent)
	txn := db.DB.Where("id = ?", id).Preload("NotificationCampaign").Preload("NotificationCampaign.NotificationConfiguration").First(notificationSent)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting notification sent: %v", txn.Error)
		return notificationSent, txn.Error
	}
	return notificationSent, nil
}

func GetNotificationCampaignById(id string) (*models.NotificationCampaign, error) {
	notificationCampaign := new(models.NotificationCampaign)
	txn := db.DB.Where("id = ?", id).Preload("NotificationConfiguration").First(notificationCampaign)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting notification campaign: %v", txn.Error)
		return notificationCampaign, txn.Error
	}
	return notificationCampaign, nil
}

func GetShopFromShopIdentifier(shopIdentifier string) (*models.Shop, error) {
	shop := new(models.Shop)
	txn := db.DB.Where("shop_identifier = ?", shopIdentifier).First(shop)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting shop: %v", txn.Error)
		return shop, txn.Error
	}
	return shop, nil
}

func GetNotificationsSentByCampaignId(campaignId string) ([]models.NotificationsSent, error) {
	notificationsSent := []models.NotificationsSent{}
	txn := db.DB.Where("notification_campaign_id = ?", campaignId).Preload("NotificationCampaign").Find(&notificationsSent)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting notifications sent: %v", txn.Error)
		return notificationsSent, txn.Error
	}
	return notificationsSent, nil
}

func GetNotificationCampaignsByShopId(shopId string, notification_campaign_id string) ([]models.NotificationCampaign, error) {
	notificationCampaigns := []models.NotificationCampaign{}
	var txn *gorm.DB
	if notification_campaign_id != "" {
		txn = db.DB.Where("shop_id = ? AND id = ?", shopId, notification_campaign_id).Preload("Shop").Preload("NotificationConfiguration").Find(&notificationCampaigns)
	} else {
		txn = db.DB.Where("shop_id = ?", shopId).Preload("Shop").Preload("NotificationConfiguration").Find(&notificationCampaigns)
	}

	if txn.Error != nil {
		log.Printf("[ERROR] Error getting campaigns: %v", txn.Error)
		return notificationCampaigns, txn.Error
	}
	return notificationCampaigns, nil
}

func GetStoreNotifications(shopId, notificationType string) ([]models.Notification, error) {	
	notifications := []models.Notification{}
	query := db.DB.Where("shop_id = ? AND notification_type = ?", shopId, notificationType)

	if notificationType == "*" {
		query = db.DB.Where("shop_id = ?", shopId)
	}

	txn := query.Find(&notifications)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting notifications: %v", txn.Error)
		return notifications, txn.Error
	}
	return notifications, nil
}

func SetTrackedClick(trackedClick *models.TrackedClick) (*models.TrackedClick, error) {
	// check if tracked click with ID exists
	if trackedClick.ID == "" {
		trackedClick.CreatedAt = db.DB.NowFunc().String()
		trackedClick.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("NotificationCampaign").Create(trackedClick)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating tracked click: %v", txn.Error)
			return trackedClick, txn.Error
		}
	} else {
		trackedClick.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("NotificationCampaign").Save(trackedClick)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving tracked click: %v", txn.Error)
			return trackedClick, txn.Error
		}
	}

	return trackedClick, nil
}

func SetNotification(notification *models.Notification) (*models.Notification, error) {
	// check if notification with ID exists
	if notification.ID == "" {
		notification.CreatedAt = db.DB.NowFunc().String()
		notification.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Create(notification)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating notification: %v", txn.Error)
			return notification, txn.Error
		}
	} else {
		notification.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Save(notification)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving notification: %v", txn.Error)
			return notification, txn.Error
		}
	}

	return notification, nil
}

func GetTrackedUserByFingerprint(fingerprint string) (*models.TrackedUser, error) {
	trackedUser := new(models.TrackedUser)
	txn := db.DB.Where("fingerprint = ?", fingerprint).First(&trackedUser)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting tracked user: %v", txn.Error)
		return trackedUser, txn.Error
	}
	return trackedUser, nil
}

func AppendCustomerIDToSubscription(subscription *models.NotificationSubscription, customerID string) (error) {
	// remember, subscription.CustomerIDs is a pg.Int64Array
	return db.DB.Model(
		&subscription,
	).Where(
		"id = ?",
		subscription.ID,
	).Omit("Shop").Update(
		"customer_ids",
		gorm.Expr("array_append(customer_ids, ?)", customerID),
	).Error
}

func SetTrackedUser(trackedUser *models.TrackedUser) (*models.TrackedUser, error) {
	// check if tracked user with ID exists
	if trackedUser.ID == "" {
		// check if user with fingerprint exists
		trackedUser.CreatedAt = db.DB.NowFunc().String()
		trackedUser.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Create(trackedUser)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating tracked user: %v", txn.Error)
			return trackedUser, txn.Error
		}
	} else {
		trackedUser.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Save(trackedUser)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving tracked user: %v", txn.Error)
			return trackedUser, txn.Error
		}
	}

	return trackedUser, nil
}

func SetUser(user *models.User) (*models.User, error) {
	// check if user with ID exists
	if user.ID == "" {
		user.CreatedAt = db.DB.NowFunc().String()
		user.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("CurrentShop").Create(user)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating user: %v", txn.Error)
			return user, txn.Error
		}
	} else {
		user.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("CurrentShop").Save(user)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving user: %v", txn.Error)
			return user, txn.Error
		}
	}

	return user, nil
}

func SetNotficationSubscription(subscription models.NotificationSubscription) (models.NotificationSubscription, error) {
	// check if subscription with ID exists
	if subscription.ID == "" {
		subscription.ID = uuid.New().String()
		subscription.CreatedAt = db.DB.NowFunc().String()
		subscription.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("Shop").Create(&subscription)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating subscription: %v", txn.Error)
			return subscription, txn.Error
		}
	} else {
		subscription.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("Shop").Save(&subscription)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving subscription: %v", txn.Error)
			return subscription, txn.Error
		}
	}

	return subscription, nil
}

func SetNotificationCampaign(notificationCampaign *models.NotificationCampaign) (models.NotificationCampaign, error) {
	// check if notificationCampaign with ID exists
	if notificationCampaign.ID == "" {
		notificationCampaign.ID = uuid.New().String()
		notificationCampaign.CreatedAt = db.DB.NowFunc().String()
		notificationCampaign.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("NotificationConfiguration", "Shop").Create(notificationCampaign)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating notification campaign: %v", txn.Error)
			return *notificationCampaign, txn.Error
		}
	} else {
		notificationCampaign.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("NotificationConfiguration", "Shop").Save(notificationCampaign)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving notification campaign: %v", txn.Error)
			return *notificationCampaign, txn.Error
		}
	}

	return *notificationCampaign, nil
}

func SetNotificationSent(notificationSent *models.NotificationsSent) (models.NotificationsSent, error) {
	if notificationSent.ID == "" {
		notificationSent.CreatedAt = db.DB.NowFunc().String()
		notificationSent.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Omit("NotificationCampaign").Create(notificationSent)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating notification sent: %v", txn.Error)
			return *notificationSent, txn.Error
		}
	} 

	notificationSent.UpdatedAt = db.DB.NowFunc().String()
	txn := db.DB.Omit("NotificationCampaign").Save(notificationSent)
	if txn.Error != nil {
		log.Printf("[ERROR] Error saving notification sent: %v", txn.Error)
		return *notificationSent, txn.Error
	}

	return *notificationSent, nil
}

func DeleteAllUserOwnedNotificationSubscriptions(ownerID string) error {
	// Ensure the model is a pointer
	var notificationSubscription models.NotificationSubscription

	// Enable GORM debug mode for detailed logs
	db.DB = db.DB.Debug()

	// Perform the delete operation
	txn := db.DB.Where("owner_id = ?", ownerID).Delete(&notificationSubscription)
	if txn.Error != nil {
		log.Printf("[ERROR] Error deleting notification subscriptions: %v", txn.Error)
		return txn.Error
	}

	// Log the number of rows affected
	log.Printf("[INFO] Rows affected: %d", txn.RowsAffected)
	return nil
}

func DeleteNotification(notification *models.Notification) error {
	// Enable GORM debug mode for detailed logs
	db.DB = db.DB.Debug()

	// Perform the delete operation
	txn := db.DB.Delete(notification)
	if txn.Error != nil {
		log.Printf("[ERROR] Error deleting notification: %v", txn.Error)
		return txn.Error
	}

	// Log the number of rows affected
	log.Printf("[INFO] Rows affected: %d", txn.RowsAffected)
	return nil
}

