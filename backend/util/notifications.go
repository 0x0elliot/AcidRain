package util

import (
	"encoding/json"
	"log"
	"os"
	"io/ioutil"

	webpush "github.com/SherClockHolmes/webpush-go"

	"go-authentication-boilerplate/models"
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
func SendPushNotification(title string, message string, icon string, badge string, url string, subscriptionId string, notificationCampaign models.NotificationCampaign) error {	
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

	if notificationCampaign.ID != "" {
		// localhost:5002/api/tracking/redirect/click?notifications_sent_id=5e34818e-8606-44d8-bf4b-e25b7326b431
		url = os.Getenv("ACIDRAIN_API_URL") + "/api/tracking/redirect/click?notification_campaign_id=" + notificationCampaign.ID
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
	var notificationSent models.NotificationsSent

	if notificationCampaign.ID != "" {
		// record the notification
		notificationSent = models.NotificationsSent{
			NotificationCampaignID: notificationCampaign.ID,
			Status: "pending",
		}
	}

	notif, errWebPush := webpush.SendNotification([]byte(totalMessage), &pushSubscription, &webpush.Options{
		Subscriber:      "",
		VAPIDPublicKey: publicKey,
		VAPIDPrivateKey: privateKey,
		TTL:             30,
	})

	log.Printf("[INFO] Notification: %v", notif)


	if notificationCampaign.ID != "" {
		type Meta struct {
			Body string `json:"body"`
			Headers map[string][]string `json:"headers"`
		}

		var meta Meta

		body, err := ioutil.ReadAll(notif.Body)
		if err != nil {
			log.Printf("[ERROR] Error reading notification body: %v", err)
			meta.Body = "COULD NOT READ BODY -- " + err.Error()
		} else {
			meta.Body = string(body)
		}

		meta.Headers = notif.Header
		
		// now, convert all this to a JSON string
		metaJson, err := json.Marshal(meta)
		if err != nil {
			log.Printf("[ERROR] Error marshalling meta: %v", err)
			return err
		}

		notificationSent.APIResponse = string(metaJson)

		// now, get all the headers, and store them too
	
		notificationSent.APIStatus = notif.StatusCode
		notificationSent, err = SetNotificationSent(&notificationSent)
		if err != nil {
			log.Printf("[ERROR] Error setting notification sent: %v", err)
		}
	}

	if errWebPush != nil {
		log.Printf("[ERROR] Error sending push notification: %v", err)

		if notificationCampaign.ID != "" {
			notificationSent.Status = "failed"
		}

		return errWebPush
	}

	log.Printf("[INFO] Push notification sent successfully")
	if notificationCampaign.ID != "" {
		notificationSent.Status = "sent"
		_, err = SetNotificationSent(&notificationSent)
		if err != nil {
			log.Printf("[ERROR] Error setting notification sent: %v", err)
		}
	}

	return nil
}
