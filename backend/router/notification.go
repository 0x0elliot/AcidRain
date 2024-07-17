package router

import (
	util "go-authentication-boilerplate/util"
	auth "go-authentication-boilerplate/auth"
	models "go-authentication-boilerplate/models"

	"log"
	// "strings"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes() {
	NOTIFICATION.Post("/subscribe", HandlePublicSubscribeToPush)

	// set up
	privNotification := NOTIFICATION.Group("/private")

	privNotification.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privNotification.Post("/push", HandlePushNotification)
	privNotification.Post("/subscribe", HandleSubscribeToPush)
}

type SubscribeToPushRequest struct {
	Endpoint string `json:"endpoint"`
	ExpirationTime int64 `json:"expirationTime"`
	Keys struct {
		P256dh string `json:"p256dh"`
		Auth string `json:"auth"`
	} `json:"keys"`
}

func HandlePublicSubscribeToPush(c *fiber.Ctx) error {
	type PublicSubscribeToPushRequest struct {
		Subscription SubscribeToPushRequest `json:"subscription"`
		StoreUrl string `json:"storeUrl"`
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
	subscription.OwnerID = shop.OwnerID

	// subscribe user to push
	err = util.SubscribeUserToPush(subscription, shop.OwnerID)
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
	err := util.SubscribeUserToPush(subscription, c.Locals("id").(string))
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


