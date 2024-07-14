package router

import (
	util "go-authentication-boilerplate/util"
	auth "go-authentication-boilerplate/auth"
	models "go-authentication-boilerplate/models"

	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes() {
	// set up
	privNotification := NOTIFICATION.Group("/private")

	privNotification.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privNotification.Post("/push", HandlePushNotification)
	privNotification.Post("/subscribe", HandleSubscribeToPush)
}

func HandleSubscribeToPush(c *fiber.Ctx) error {
	type SubscribeToPushRequest struct {
		Endpoint string `json:"endpoint"`
		ExpirationTime int64 `json:"expirationTime"`
		Keys struct {
			P256dh string `json:"p256dh"`
			Auth string `json:"auth"`
		} `json:"keys"`
	}

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

		for _, subscription := range subscriptions {
			err := util.SendPushNotification(req.Body, subscription.ID)
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

	// send push notification
	err := util.SendPushNotification(req.Body, c.Locals("id").(string))
	if err != nil {
		log.Printf("[ERROR] Error in sending push notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error sending push notification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Push notification sent successfully",
	})
}


