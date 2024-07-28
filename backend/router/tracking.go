package router

import (
	models "go-authentication-boilerplate/models"
	util "go-authentication-boilerplate/util"

	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupTrackingRoutes() {
	// set up
	// TRACKING.Post("/sync", HandleTrackedUserSync)
	TRACKING.Get("/redirect/click", HandleTrackClick)
}

func HandleTrackClick(c *fiber.Ctx) error {
	// this is just a forward function. it is supposed
	// to be as fast as possible. the query will have the notifications_sent_id
	// just forward user to notification_campaign.URL, save the click info as well

	// get the notification_campaign_id
	notification_sent_id := c.Query("notifications_sent_id")
	notificationSent, err := util.GetNotificationSentById(notification_sent_id)
	if err != nil {
		log.Printf("[ERROR] Error getting notification sent: %v -- %v", notification_sent_id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error getting notification sent",
		})
	}

	notificationCampaign := notificationSent.NotificationCampaign


	headers := c.Request().Header

	// save the click
	click := new(models.TrackedClick)
	click.NotificationCampaignID = notificationCampaign.ID
	click.NotificationCampaign = *notificationCampaign
	click.ClickHeaders = headers.String()
	click.ClickIP = c.IP()
	click.ClickForwardURL = notificationCampaign.NotificationConfiguration.URL
	click.ClickReferrer = c.Get("Referer")
	click.ClickUserAgent = c.Get("User-Agent")
	click.ClickOS = c.Get("OS")

	_, err = util.SetTrackedClick(click)
	if err != nil {
		log.Printf("[ERROR] Error saving click: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error saving click",
		})
	}

	// redirect to the URL
	return c.Redirect(notificationCampaign.NotificationConfiguration.URL, fiber.StatusSeeOther)
}

// alright, GDPR will try to fuck me in the ass. 
// let's not go this route right now.
// add as less validation here as possible by design
// be open to all blocking edge cases. This is just a plain WRITE APIs and if else conditions.
// that's how they rule the world. We will too.
// "Do not let the hero in your soul perish, in lonely frustration for
// the life you deserved, but have never been able to reach. The world you desire
// can be won. It exists.. It is real.. it is possible..
// It's yours." - Ayn Rand"
// func HandleTrackedUserSync(c *fiber.Ctx) error {
// 	type FingerPrint struct {
// 		FingerPrint                  string `json:"fingerprint"`
// 		ShopifyUniqueID              string `json:"shopify_unique_id"`
// 		Store                        string `json:"store"`
// 		PushNotificationSubscription string `json:"push_notification_subscription"`
// 	}

// 	newFingerPrint := new(FingerPrint)
// 	if err := c.BodyParser(newFingerPrint); err != nil {
// 		log.Printf("[ERROR] Error parsing request body: %v", err)
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Error parsing request body",
// 		})
// 	}

// 	trackedUser := new(models.TrackedUser)
// 	trackedUser.FingerPrint = newFingerPrint.FingerPrint

// 	// adjust the following to the use of pq.StringArray
// 	trackedUser.ShopifyUniqueIDs = append(trackedUser.ShopifyUniqueIDs, newFingerPrint.ShopifyUniqueID)

// 	if newFingerPrint.PushNotificationSubscription != "" {
// 		trackedUser.SubscriptionMetadata = append(trackedUser.SubscriptionMetadata, newFingerPrint.PushNotificationSubscription)
// 	}
// 	trackedUser.Stores = append(trackedUser.Stores, newFingerPrint.Store)

// 	trackedUserOriginal, err := util.GetTrackedUserByFingerprint(newFingerPrint.FingerPrint)
// 	if err != nil {
// 		log.Printf("[ERROR] Error getting fingerprint: %v", err)
// 	} else {
// 		trackedUser.ShopifyUniqueIDs = append(trackedUser.ShopifyUniqueIDs, trackedUserOriginal.ShopifyUniqueIDs...)
// 		if newFingerPrint.PushNotificationSubscription != "" {
// 			trackedUser.SubscriptionMetadata = append(trackedUser.SubscriptionMetadata, trackedUserOriginal.SubscriptionMetadata...)
// 		}
// 		trackedUser.Stores = append(trackedUser.Stores, trackedUserOriginal.Stores...)

// 		trackedUserOriginal.ShopifyUniqueIDs = trackedUser.ShopifyUniqueIDs
// 		if newFingerPrint.PushNotificationSubscription != "" {
// 			trackedUserOriginal.SubscriptionMetadata = trackedUser.SubscriptionMetadata
// 		}
// 		trackedUserOriginal.Stores = trackedUser.Stores

// 		trackedUser = trackedUserOriginal
// 	}

// 	// save the fingerprint
// 	userId, err := util.SetTrackedUser(trackedUser)
// 	if err != nil {
// 		log.Printf("[ERROR] Error in saving fingerprint: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error":   true,
// 			"message": "Error saving fingerprint",
// 			"user_id": userId,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"error":   false,
// 		"message": "Fingerprint saved successfully",
// 	})
// }
