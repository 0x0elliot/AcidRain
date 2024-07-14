package util

import (
	"go-authentication-boilerplate/models"
	"log"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
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
	publicKey := os.Getenv("ACIDRAIN_WEB_PUSH_PUBLIC_KEY")
	privateKey := os.Getenv("ACIDRAIN_WEB_PUSH_PRIVATE_KEY")

	subscription, err := GetNotificationSubscriptionById(subscriptionId)
	if err != nil {
		log.Printf("[ERROR] Error getting subscription: %v", err)
		return err
	}

	pushSubscription := webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.Auth,
			P256dh: subscription.P256dh,
		},
	}

	_, err = webpush.SendNotification([]byte(message), &pushSubscription, &webpush.Options{
		Subscriber:      "You're subscribed to AcidRain",
		VAPIDPublicKey: publicKey,
		VAPIDPrivateKey: privateKey,
		TTL:             30,
	}) 

	if err != nil {
		log.Printf("[ERROR] Error sending push notification: %v", err)
		return err
	}

	log.Printf("[INFO] Push notification sent successfully")

	return nil
}
