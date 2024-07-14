package util

import (
	"go-authentication-boilerplate/models"
	"log"

	// webpush "github.com/SherClockHolmes/webpush-go"
)

func SubscribeUserToPush(subscription models.NotificationSubscription, userId string) error {
	subscription.ID = ""

	if userId == "" {
		// API is in actual use
	} else {
		// user is trying to subscribe to test
		log.Printf("[INFO] User %v is trying to subscribe to push notifications", userId)
		_, err := GetUserById(userId)
		if err != nil {
			log.Printf("[ERROR] Error getting user: %v", err)
			return err
		}

		err = DeleteAllUserOwnedNotificationSubscriptions(userId)
		if err != nil {
			log.Printf("[ERROR] Error deleting user's subscriptions: %v", err)
		}

		subscription.OwnerID = userId
	}

	sub, err := SetNotficationSubscription(subscription)
	if err != nil {
		log.Printf("[ERROR] Error setting subscription: %v", err)
		return err
	}

	log.Printf("[INFO] Subscription set: %v", sub)
	return nil
}

func SendPushNotification(message string, subscriptionId string) error {
	// API is in actual use
	return nil
}
