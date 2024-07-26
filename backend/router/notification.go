package router

import (
	"fmt"
	auth "go-authentication-boilerplate/auth"
	models "go-authentication-boilerplate/models"
	util "go-authentication-boilerplate/util"

	"log"
	// "strings"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes() {
	NOTIFICATION.Post("/subscribe", HandlePublicSubscribeToPush)
	NOTIFICATION.Post("/sync", HandlePublicSync)

	// set up
	privNotification := NOTIFICATION.Group("/private")

	privNotification.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privNotification.Post("/push", HandlePushNotification)
	privNotification.Post("/subscribe", HandleSubscribeToPush)
	privNotification.Get("/notifications", HandleGetNotifications)

	privNotification.Post("/enable/push-notifications", HandleEnablePushNotifications)
	privNotification.Post("/disable/push-notifications", HandleDisablePushNotifications)

	privNotification.Get("/push-subscribers", HandleGetPushSubscribers)
}

type SubscribeToPushRequest struct {
	Endpoint string `json:"endpoint"`
	ExpirationTime int64 `json:"expirationTime"`
	Keys struct {
		P256dh string `json:"p256dh"`
		Auth string `json:"auth"`
	} `json:"keys"`
}

func HandlePublicSync(c *fiber.Ctx) error {
	type PublicSyncRequest struct {
		Subscription SubscribeToPushRequest `json:"subscription"`
		StoreUrl string `json:"storeUrl"`
		Cid int64 `json:"cid"`
	}

	var req PublicSyncRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	shop, err := util.GetShopFromShopIdentifier(req.StoreUrl)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	// get existing subscription
	subscription, err := util.GetSubscriptionFromEndpoint(req.Subscription.Endpoint)
	if err != nil {
		log.Printf("[ERROR] Error getting subscription: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting subscription",
		})
	}

	// if subscription.OwnerID != shop.OwnerID {
	if subscription.Shop.ID != shop.ID {
		log.Printf("[ERROR] Unauthorized access")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message":   "Unauthorized access",
		})
	}

	// update subscription
	subscription.Auth = req.Subscription.Keys.Auth
	subscription.P256dh = req.Subscription.Keys.P256dh

	sub, err := util.SetNotficationSubscription(*subscription)
	if err != nil {
		log.Printf("[ERROR] Error setting subscription: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error setting subscription",
		})
	}

	existingCustomerIds := sub.CustomerIDs
	if !util.Contains(existingCustomerIds, fmt.Sprint(req.Cid)) {
		err := util.AppendCustomerIDToSubscription(&sub, fmt.Sprint(req.Cid))
		if err != nil {
			log.Printf("[ERROR] Error appending customer id to subscription: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message":   "Error appending customer id to subscription",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Subscription updated successfully",
	})
}

func HandleGetNotifications(c *fiber.Ctx) error {
	type GetNotificationsRequest struct {
		ShopIdentifier string `json:"shop_identifier"`
	}

	var req GetNotificationsRequest
	// this is a GET request, so we need to get the query params
	req.ShopIdentifier = c.Query("shop_identifier")

	shop, err := util.GetShopFromShopIdentifier(req.ShopIdentifier)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Unauthorized access")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message":   "Unauthorized access",
		})
	}

	notifications, err := util.GetStoreNotifications(shop.ID, "*")
	if err != nil {
		log.Printf("[ERROR] Error getting notifications: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting notifications",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"notifications": notifications,
	})
}

func HandleGetPushSubscribers(c *fiber.Ctx) error {
	type GetPushSubscribersRequest struct {
		ShopIdentifier string `json:"shop_identifier"`
	}

	var req GetPushSubscribersRequest
	// this is a GET request, so we need to get the query params
	req.ShopIdentifier = c.Query("shop_identifier")

	shop, err := util.GetShopFromShopIdentifier(req.ShopIdentifier)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Unauthorized access")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message":   "Unauthorized access",
		})
	}

	subscriptions, err := util.GetNotificationSubscriptionByShopId(shop.ID)
	if err != nil {
		log.Printf("[ERROR] Error getting subscriptions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting subscriptions",
		})
	}

	for i, _ := range subscriptions {
		subscriptions[i].Auth = ""
		subscriptions[i].P256dh = ""
		subscriptions[i].Endpoint = ""

		log.Printf("[DEBUG] Shop: %v", subscriptions[i].Shop)

		subscriptions[i].Shop.AccessToken = ""
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"subscriptions": subscriptions,
	})
}

func HandleDisablePushNotifications(c *fiber.Ctx) error {
	type DisablePushNotificationsRequest struct {
		ShopIdentifier string `json:"shop_identifier"`
	}

	var req DisablePushNotificationsRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	shop, err := util.GetShopFromShopIdentifier(req.ShopIdentifier)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Unauthorized access")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message":   "Unauthorized access",
		})
	}

	notifs, err := util.GetStoreNotifications(shop.ID, "push")
	if err != nil {
		log.Printf("[ERROR] Error getting notifications: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting notifications",
		})
	}

	if len(notifs) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"message":   "Push notifications already disabled for this store",
		})
	}

	for _, notif := range notifs {
		err = util.DeleteNotification(&notif)
		if err != nil {
			log.Printf("[ERROR] Error deleting notification: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message":   "Error deleting notification",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Push notifications disabled successfully",
	})
}
func HandleEnablePushNotifications(c *fiber.Ctx) error {
	type EnablePushNotificationsRequest struct {
		ShopIdentifier string `json:"shop_identifier"`
	}

	var req EnablePushNotificationsRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	shop, err := util.GetShopFromShopIdentifier(req.ShopIdentifier)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Unauthorized access")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message":   "Unauthorized access",
		})
	}

	notifs, err := util.GetStoreNotifications(shop.ID, "push")
	if err != nil {
		log.Printf("[ERROR] Error getting notifications: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting notifications",
		})
	}

	if len(notifs) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"message":   "Push notifications already enabled for this store",
		})
	}

	var notification models.Notification
	notification.ShopID = shop.ID
	notification.NotificationType = "push"

	// configuration includes the message, title, image, etc
	notification.Configured = false

	_, err = util.SetNotification(&notification)
	if err != nil {
		log.Printf("[ERROR] Error setting notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error setting notification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Push notifications enabled successfully",
	})
}


func HandlePublicSubscribeToPush(c *fiber.Ctx) error {
	type PublicSubscribeToPushRequest struct {
		Subscription SubscribeToPushRequest `json:"subscription"`
		StoreUrl string `json:"storeUrl"`
		Customer struct {
			Cid int64 `json:"cid"`
		} `json:"customer"`
	}

	var req PublicSubscribeToPushRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	shop, err := util.GetShopFromShopIdentifier(req.StoreUrl)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	var subscription models.NotificationSubscription
	subscription.Endpoint = req.Subscription.Endpoint
	subscription.Auth = req.Subscription.Keys.Auth
	subscription.P256dh = req.Subscription.Keys.P256dh
	// subscription.OwnerID = shop.OwnerID
	// subscription.ShopID = shop.ID
	subscription.Shop = *shop

	// check if subscription already exists
	_, err = util.GetIdenticalSubscription(subscription)
	if err != nil {
		if err.Error() == "record not found" {
			// normal flow
		} else {
			log.Printf("[ERROR] Error getting subscription: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message":   "Error getting subscription",
			})
		}
	} else {
		log.Printf("[INFO] Subscription already exists")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"message":   "Subscription already exists",
		})
	}

	// subscribe user to push
	sub, err := util.SubscribeUserToPush(subscription, "")
	if err != nil {
		log.Printf("[ERROR] Error in subscribing user to push: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error subscribing user to push",
		})
	} else {
		customerCidStr := fmt.Sprint(req.Customer.Cid)
		if customerCidStr != "0" {
			err := util.AppendCustomerIDToSubscription(&sub, customerCidStr)
			if err != nil {
				log.Printf("[ERROR] Error appending customer id to subscription: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"message":   "Error appending customer id to subscription",
				})
			}
		}	
	}

	go func() {
		// send push notification
		err := util.SendPushNotification("Welcome to AcidRain", "You're now subscribed to push notifications", "https://upload.wikimedia.org/wikipedia/en/a/a9/Example.jpg", "https://upload.wikimedia.org/wikipedia/en/a/a9/Example.jpg", "http://localhost:3000", sub.ID)
		if err != nil {
			log.Printf("[ERROR] Error in sending push notification: %v", err)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "User subscribed to push notifications successfully",
	})
}

func HandleSubscribeToPush(c *fiber.Ctx) error {
	var req SubscribeToPushRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	var subscription models.NotificationSubscription
	subscription.Endpoint = req.Endpoint
	subscription.Auth = req.Keys.Auth
	subscription.P256dh = req.Keys.P256dh

	// subscribe user to push
	_, err := util.SubscribeUserToPush(subscription, c.Locals("id").(string))
	if err != nil {
		log.Printf("[ERROR] Error in subscribing user to push: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error subscribing user to push",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "User subscribed to push notifications successfully",
	})
}

func HandlePushNotification(c *fiber.Ctx) error {
	type PushNotificationRequest struct {
		Test bool `json:"test"`
		Body string `json:"body"`
		Title string `json:"title"`
	}

	var req PushNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	if req.Test {
		log.Printf("[INFO] Test push notification requested")
		subscriptions, err := util.GetNoficationSubscriptionByOwnerId(c.Locals("id").(string))
		if err != nil {
			log.Printf("[ERROR] Error getting subscriptions: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message":   "Error getting subscriptions",
			})
		}

		if len(subscriptions) == 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": false,
				"message":   "No subscriptions found",
			})
		}

		examplePic := "https://upload.wikimedia.org/wikipedia/en/a/a9/Example.jpg"

		for _, subscription := range subscriptions {
			req.Body = "Test push notification"
			req.Title = "You're now subscribed!"
			// err := util.SendPushNotification(req.Title, req.Body, subscription.ID)
			err := util.SendPushNotification(req.Title, req.Body, examplePic, examplePic,"http://localhost:3000", subscription.ID) 
			if err != nil {
				log.Printf("[ERROR] Error in sending push notification: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"message":   "Error sending push notification",
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"message":   "Push notification sent successfully",
		})
	}

	// // send push notification
	// err := util.SendPushNotification(req.Title, req.Body, c.Locals("id").(string))
	// if err != nil {
	// 	log.Printf("[ERROR] Error in sending push notification: %v", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"message":   "Error sending push notification",
	// 	})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Push notification sent successfully",
	})
}


