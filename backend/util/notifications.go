package util

import (
	"encoding/json"
	"go-authentication-boilerplate/models"
	"log"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// the assumption here is that if this function is called, this is the first subscription
// for the user
func SubscribeUserToPush(subscription models.NotificationSubscription, userId string) (models.NotificationSubscription, error) {
	var sub models.NotificationSubscription

	// check if subscription with this endpoint already exists
	subPtr, err := GetIdenticalSubscription(subscription)
	if err != nil {
		if err.Error() == "record not found" {
		} else {
			log.Printf("[ERROR] Error getting identical subscription: %v", err)
			return sub, err
		}
	}

	sub = *subPtr

	if sub.ID != "" {
		log.Printf("[INFO] Subscription already exists: %v", sub)
		return sub, nil
	}

	subscription.ID = ""

	if userId == "" {
		// actual user is trying to subscribe
	} else {
		// user is trying to subscribe to test
		log.Printf("[INFO] User %v is trying to subscribe to push notifications", userId)
		user_, err := GetUserById(userId)
		if err != nil {
			log.Printf("[ERROR] Error getting user: %v", err)
			return sub, err
		}

		err = DeleteAllUserOwnedNotificationSubscriptions(userId)
		if err != nil {
			log.Printf("[ERROR] Error deleting user's subscriptions: %v", err)
		}

		subscription.OwnerID = userId
		subscription.ShopID = *user_.CurrentShopID
		subscription.Shop = user_.CurrentShop
	}

	log.Printf("[INFO] Subscription: %v", subscription)

	sub, err = SetNotficationSubscription(subscription)
	if err != nil {
		log.Printf("[ERROR] Error setting subscription: %v", err)
		return sub, err
	}

	log.Printf("[INFO] Subscription set: %v", sub)
	return sub, nil
}


// ideally, all "url" should be an analytics redirect URL from our end
func SendPushNotification(title string, message string, icon string, badge string, url string, subscriptionId string) error {
	type PushNotificationRequest struct {
		Body string `json:"body"`
		Title string `json:"title"`
		Icon string `json:"icon"`
		URL string `json:"url"`
		Badge string `json:"badge"`
	}

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

	totalMessageStruct := PushNotificationRequest{
		Body: message,
		Title: title,
		Icon: icon,
		URL: url,
		Badge: badge,
	}

	// make sure to convert the struct to a JSON string
	jsonData, err := json.Marshal(totalMessageStruct)
	if err != nil {
		log.Printf("[ERROR] Error marshalling message: %v", err)
		return err
	}

	totalMessage := string(jsonData)

	notif, err := webpush.SendNotification([]byte(totalMessage), &pushSubscription, &webpush.Options{
		Subscriber:      "",
		VAPIDPublicKey: publicKey,
		VAPIDPrivateKey: privateKey,
		TTL:             30,
	}) 

	log.Printf("[INFO] Notification: %v", notif)

	if err != nil {
		log.Printf("[ERROR] Error sending push notification: %v", err)
		return err
	}

	log.Printf("[INFO] Push notification sent successfully")

	return nil
}
