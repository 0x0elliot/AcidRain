package router

import (
	models "go-authentication-boilerplate/models"
	util "go-authentication-boilerplate/util"
	auth "go-authentication-boilerplate/auth"

	"strings"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupTrackingRoutes() {
	// set up
	// TRACKING.Post("/sync", HandleTrackedUserSync)
	TRACKING.Get("/redirect/click", HandleTrackClick)

	privTracking := TRACKING.Group("/private")
	privTracking.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privTracking.Get("/stats/devices", HandleGetCampaignStatistics)
	privTracking.Get("/stats/os", HandleOSClickStats)
}

func HandleGetCampaignStatistics(c *fiber.Ctx) error {
	shopID := c.Query("shop_id")

	if shopID == "" {
		log.Printf("[ERROR] Shop ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Shop ID is required",
		})
	}

	shop, err := util.GetShopById(shopID)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Shop does not belong to user")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Shop does not belong to user",
		})
	}

	if c.Query("start") == "" || c.Query("end") == "" {
		log.Printf("[ERROR] Start and End dates not given. Defaulting to 1st April 2024 to 4th April 2024")
	}

	// startDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	// endDate := time.Date(2024, 8, 4, 23, 59, 59, 999999999, time.UTC)
	var startDate time.Time
	var endDate time.Time

	if c.Query("start") == "" || c.Query("end") == "" {
		startDate = time.Now().AddDate(0, -3, 0)
		endDate = time.Now()
	}


	stats, err := util.FetchClickStats(
		c.Query("shop_id"),
		startDate,
		endDate,
	)

	if err != nil {
		log.Printf("[ERROR] Error fetching stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error fetching stats",
		})
	}

	// Use the stats data as needed
	// for _, stat := range stats {
	// 	fmt.Printf("Date: %s, Desktop: %d, Mobile: %d\n", stat.Date.Format("2006-01-02"), stat.Desktop, stat.Mobile)
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"stats": stats,
	})
}

func HandleOSClickStats(c *fiber.Ctx) error {
	shopID := c.Query("shop_id")

	if shopID == "" {
		log.Printf("[ERROR] Shop ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Shop ID is required",
		})
	}

	shop, err := util.GetShopById(shopID)
	if err != nil {
		log.Printf("[ERROR] Error getting shop: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error getting shop",
		})
	}

	if shop.OwnerID != c.Locals("id").(string) {
		log.Printf("[ERROR] Shop does not belong to user")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Shop does not belong to user",
		})
	}

	var startDate time.Time
	var endDate time.Time

	if c.Query("start") == "" || c.Query("end") == "" {
		log.Printf("[ERROR] Start and End dates not given. Defaulting to 1st April 2024 to 4th April 2024")
		// startDate is 3 months ago
		startDate = time.Now().AddDate(0, -3, 0)
		endDate = time.Now()
	}

	stats, err := util.FetchOSChartData(
		c.Query("shop_id"),
		startDate,
		endDate,
	)

	if err != nil {
		log.Printf("[ERROR] Error fetching stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error fetching stats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"stats": stats,
	})
}

func HandleTrackClick(c *fiber.Ctx) error {
	// this is just a forward function. it is supposed
	// to be as fast as possible. the query will have the notifications_sent_id
	// just forward user to notification_campaign.URL, save the click info as well

	// get the notification_campaign_id
	notification_campaign_id := c.Query("notification_campaign_id")
	if notification_campaign_id == "" {
		log.Printf("[ERROR] Notification Campaign ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Notification Campaign ID is required",
		})
	}

	notificationCampaign, err := util.GetNotificationCampaignById(notification_campaign_id)
	if err != nil {
		log.Printf("[ERROR] Error getting notification campaign: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Error getting notification campaign",
		})
	}

	headers := c.Request().Header

	// save the click
	click := new(models.TrackedClick)
	click.NotificationCampaignID = notificationCampaign.ID
	click.NotificationCampaign = *notificationCampaign
	click.ClickHeaders = headers.String()
	click.ClickIP = c.IP()
	click.ClickForwardURL = notificationCampaign.NotificationConfiguration.URL
	click.ClickReferrer = c.Get("Referer")

	userAgent := c.Get("User-Agent")
	click.ClickUserAgent = userAgent

	// guess the device from the user agent
	click.ClickOS = strings.ToLower(util.GetUserAgentOS(userAgent))
	click.ClickDevice = strings.ToLower(util.GetUserAgentDeviceType(userAgent))

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
