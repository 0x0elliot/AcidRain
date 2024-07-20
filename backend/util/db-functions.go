package util

import (
	"log"
	db "go-authentication-boilerplate/database"
	models "go-authentication-boilerplate/models"
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

func GetShopFromShopIdentifier(shopIdentifier string) (*models.Shop, error) {
	shop := new(models.Shop)
	txn := db.DB.Where("shop_identifier = ?", shopIdentifier).First(shop)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting shop: %v", txn.Error)
		return shop, txn.Error
	}
	return shop, nil
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

func GetPosts(ownerID string) ([]models.Post, error) {
	posts := []models.Post{}
	txn := db.DB.Where("owner_id = ?", ownerID).Find(&posts)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting posts: %v", txn.Error)
		return posts, txn.Error
	}
	return posts, nil
}

func GetPost(id string) (*models.Post, error) {
	post := new(models.Post)
	txn := db.DB.Where("id = ?", id).First(post)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting post: %v", txn.Error)
		return post, txn.Error
	}
	return post, nil
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

func SetPost(post *models.Post) (*models.Post, error) {
	// check if post with ID exists
	if post.ID == "" {
		post.CreatedAt = db.DB.NowFunc().String()
		post.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Create(post)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating post: %v", txn.Error)
			return post, txn.Error
		}
	} else {
		post.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Save(post)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving post: %v", txn.Error)
			return post, txn.Error
		}
	}

	return post, nil
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

func SetNotficationSubscription(subscription models.NotificationSubscription) (models.NotificationSubscription, error) {
	// check if subscription with ID exists
	if subscription.ID == "" {
		subscription.CreatedAt = db.DB.NowFunc().String()
		subscription.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Create(&subscription)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating subscription: %v", txn.Error)
			return subscription, txn.Error
		}
	} else {
		subscription.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Save(&subscription)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving subscription: %v", txn.Error)
			return subscription, txn.Error
		}
	}

	return subscription, nil
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

