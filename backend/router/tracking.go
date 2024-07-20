package router

import (
	models "go-authentication-boilerplate/models"
	util "go-authentication-boilerplate/util"

	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupTrackingRoutes() {
	// set up
	TRACKING.Post("/sync", HandleTrackedUserSync)
}

// add as less validation here as possible by design
// be open to all blocking edge cases. This is just a plain WRITE APIs and if else conditions.
// that's how they rule the world. We will too.
// "Do not let the hero in your soul perish, in lonely frustration for
// the life you deserved, but have never been able to reach. The world you desire
// can be won. It exists.. It is real.. it is possible..
// It's yours." - Ayn Rand"
func HandleTrackedUserSync(c *fiber.Ctx) error {
	type FingerPrint struct {
		FingerPrint                  string `json:"fingerprint"`
		ShopifyUniqueID              string `json:"shopify_unique_id"`
		Store                        string `json:"store"`
		PushNotificationSubscription string `json:"push_notification_subscription"`
	}

	newFingerPrint := new(FingerPrint)
	if err := c.BodyParser(newFingerPrint); err != nil {
		log.Printf("[ERROR] Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error parsing request body",
		})
	}

	trackedUser := new(models.TrackedUser)
	trackedUser.FingerPrint = newFingerPrint.FingerPrint

	// adjust the following to the use of pq.StringArray
	trackedUser.ShopifyUniqueIDs = append(trackedUser.ShopifyUniqueIDs, newFingerPrint.ShopifyUniqueID)

	if newFingerPrint.PushNotificationSubscription != "" {
		trackedUser.SubscriptionMetadata = append(trackedUser.SubscriptionMetadata, newFingerPrint.PushNotificationSubscription)
	}
	trackedUser.Stores = append(trackedUser.Stores, newFingerPrint.Store)

	trackedUserOriginal, err := util.GetTrackedUserByFingerprint(newFingerPrint.FingerPrint)
	if err != nil {
		log.Printf("[ERROR] Error getting fingerprint: %v", err)
	} else {
		trackedUser.ShopifyUniqueIDs = append(trackedUser.ShopifyUniqueIDs, trackedUserOriginal.ShopifyUniqueIDs...)
		if newFingerPrint.PushNotificationSubscription != "" {
			trackedUser.SubscriptionMetadata = append(trackedUser.SubscriptionMetadata, trackedUserOriginal.SubscriptionMetadata...)
		}
		trackedUser.Stores = append(trackedUser.Stores, trackedUserOriginal.Stores...)

		trackedUserOriginal.ShopifyUniqueIDs = trackedUser.ShopifyUniqueIDs
		if newFingerPrint.PushNotificationSubscription != "" {
			trackedUserOriginal.SubscriptionMetadata = trackedUser.SubscriptionMetadata
		}
		trackedUserOriginal.Stores = trackedUser.Stores

		trackedUser = trackedUserOriginal
	}

	// save the fingerprint
	userId, err := util.SetTrackedUser(trackedUser)
	if err != nil {
		log.Printf("[ERROR] Error in saving fingerprint: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error saving fingerprint",
			"user_id": userId,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Fingerprint saved successfully",
	})
}
