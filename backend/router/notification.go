package router

import (
	"fmt"
	auth "go-authentication-boilerplate/auth"
	models "go-authentication-boilerplate/models"
	util "go-authentication-boilerplate/util"
	"strings"

	"log"
	// "strings"
	"os"

	"cloud.google.com/go/storage"
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
	privNotification.Get("/configurations", HandleGetNotificationConfigurations)

	privNotification.Post("/notification-configuration", HandleSaveNotificationConfiguration)
	privNotification.Post("/launch", HandleLaunchNotification)
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
		Customer struct {
			Cid int64 `json:"cid"`
		} `json:"customer"`
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

	log.Printf("[DEBUG Subscription shop ID and shop ID: %v, %v", subscription.ShopID, shop.ID)

	// if subscription.OwnerID != shop.OwnerID {
	if subscription.ShopID != shop.ID {
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

	// Let's play around with this after launch
	// if (req.Customer.Cid != 0) {
	// 	resp, err := util.GetCustomer(fmt.Sprint(req.Customer.Cid), shop.AccessToken, shop.ShopIdentifier)
	// 	if err != nil {
	// 		log.Printf("[ERROR] Error getting customer: %v", err)
	// 	} else {
	// 		log.Printf("[DEBUG] Customer: %v", resp)
	// 	}
	// }

	existingCustomerIds := sub.CustomerIDs
	if !util.Contains(existingCustomerIds, fmt.Sprint(req.Customer.Cid)) {
		err := util.AppendCustomerIDToSubscription(&sub, fmt.Sprint(req.Customer.Cid))
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
		CountOnly bool `json:"count_only"`
	}

	var req GetPushSubscribersRequest
	// this is a GET request, so we need to get the query params
	req.ShopIdentifier = c.Query("shop_identifier")
	req.CountOnly = c.Query("count_only") == "true"

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

	if req.CountOnly {
		subscriptionCount, err := util.GetCountOfNotificationSubscriptionsByShopId(shop.ID)
		if err != nil {
			log.Printf("[ERROR] Error getting subscription count: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message":   "Error getting subscription count",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"count": subscriptionCount,
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
	subscription.ShopID = shop.ID
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
		err := util.SendPushNotification(
			"Subscription successful!", 
			"You're now subscribed to push notifications from "+shop.Name,
			"https://raw.githubusercontent.com/zappush/zappush.github.io/master/og-image.png", // Make customizabe 
			"https://raw.githubusercontent.com/zappush/zappush.github.io/master/og-image.png", // Make customizabe
			"https://" + shop.ShopIdentifier, // Make customizabe
			sub.ID,
		)
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

func HandleGetNotificationConfigurations(c *fiber.Ctx) error {
	type GetNotificationCampaignRequest struct {
		ShopId string `json:"shop_id"`
	}

	var req GetNotificationCampaignRequest

	req.ShopId = c.Query("shop_id")

	shop, err := util.GetShopById(req.ShopId)
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

	notifConfigs, err := util.GetNotificationConfigurationsById(shop.ID)
	if err != nil {
		log.Printf("[ERROR] Error getting notification configuration: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting notification configuration",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"configurations": notifConfigs,
	})
}


func HandleSaveNotificationConfiguration(c *fiber.Ctx) error {
	var req models.NotificationConfiguration
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	req.ID = ""
	req.CreatedAt = ""
	req.UpdatedAt = ""

	if req.ShopID == "" {
		log.Printf("[ERROR] Shop ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Shop ID is required",
		})
	}

	shop, err := util.GetShopById(req.ShopID)
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

	if (req.URL == "" || strings.HasPrefix(req.URL, "http://") || strings.HasPrefix(req.URL, "https://")) {
		req.URL = "https://" + shop.ShopIdentifier
	}

	// if configuration Icon and Badge are not empty AND they are base64 encoded
	if req.Icon != "" && req.Badge != "" {
		var storage *storage.Client
		bucketName := os.Getenv("ACIDRAIN_GCP_BUCKET_NAME")

		if util.IsBase64Image(req.Icon) || util.IsBase64Image(req.Badge) {
			storage, err = util.InitializeGCP(
				os.Getenv("ACIDRAIN_GCP_PROJECT_ID"),
				bucketName,
				os.Getenv("ACIDRAIN_GCP_CREDS"),
			)

			if err != nil {
				log.Printf("[ERROR] Error initializing GCP: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"message":   "Internal server error",
				})
			}
		}

		// check if Icon and Badge are base64 encoded
		if util.IsBase64Image(req.Icon) {
			req.Icon, err = util.UploadImageToGCP(storage, bucketName, shop.ID + "_" + req.ID + "_icon", req.Icon)
			if err != nil {
				log.Printf("[ERROR] Error uploading icon to GCP: %v", err)
				if err.Error() == "image size exceeds 5 MB" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": true,
						"message":   "Image size exceeds 5 MB",
					})
				}

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"message":   "Internal server error",
				})
			}
		}

		if util.IsBase64Image(req.Badge) {
			req.Badge, err = util.UploadImageToGCP(storage, bucketName, shop.ID + "_" + req.ID + "_badge", req.Badge)
			if err != nil {
				log.Printf("[ERROR] Error uploading badge to GCP: %v", err)
				if err.Error() == "image size exceeds 5 MB" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": true,
						"message":   "Image size exceeds 5 MB",
					})
				}

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"message":   "Error uploading badge to bucket -- is it a valid base64 image?",
				})
			}
		}
	}

	// save notification configuration
	notConfig, err := util.SetNotificationConfiguration(&req)
	if err != nil {
		log.Printf("[ERROR] Error setting notification configuration: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error setting notification configuration",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"notification": notConfig,
	})
}

func HandleLaunchNotification(c *fiber.Ctx) error {
	type LaunchNotificationRequest struct {
		ShopId string `json:"shop_id"`
		All bool `json:"all"`
		NotificationConfigurationID string `json:"notification_configuration_id"`
	}

	var req LaunchNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Error in parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Error parsing request body",
		})
	}

	if req.ShopId == "" || req.NotificationConfigurationID == "" || !req.All {
		log.Printf("[ERROR] Shop ID and Notification Configuration ID are required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Shop ID and Notification Configuration ID are required",
		})
	}

	shop, err := util.GetShopById(req.ShopId)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Forbidden request")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"message":   "Forbidden request",
		})
	}

	// get all subscriptions
	subscriptions, err := util.GetNotificationSubscriptionByShopId(shop.ID)
	if err != nil {
		log.Printf("[ERROR] Error getting subscriptions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting subscriptions",
		})
	}

	config, err := util.GetNotificationConfigurationById(req.NotificationConfigurationID, shop.ID)
	if err != nil {
		log.Printf("[ERROR] Error getting notification configuration: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting notification configuration",
		})
	}

	if len(subscriptions) == 0 {
		log.Printf("[ERROR] No subscriptions found")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"message":   "No subscriptions found",
		})
	}

	// send push notification
	for _, subscription := range subscriptions {
		err := util.SendPushNotification(
			config.Title,
			config.Message,
			config.Icon,
			config.Badge,
			config.URL,
			subscription.ID,
		)
		if err != nil {
			log.Printf("[ERROR] Error in sending push notification: %v", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Push notification sent successfully",
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

		examplePic := "https://raw.githubusercontent.com/zappush/zappush.github.io/master/og-image.png"

		for _, subscription := range subscriptions {
			req.Body = "Test push notification"
			req.Title = "Test push notification!"
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


