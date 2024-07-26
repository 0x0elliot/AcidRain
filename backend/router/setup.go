package router

import (
	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
	logger "github.com/gofiber/fiber/v2/middleware/logger"
)

var USER fiber.Router
var POST fiber.Router
var SHOP fiber.Router
var NOTIFICATION fiber.Router
var TRACKING fiber.Router

func webPushPublicKey() string {
	return "BCv7WgVIIGsZfgamKaruQEach2j6a8Us5en7Y2FIuC7PUt9aQxd2Nl2d5XIj80cfgs37DA6OE3TS1GOebJs0UTo"
}

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Change this to the allowed origins, e.g., "http://example.com"
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	// this is just for testing with shopify
	// remember that in production we will be using app proxies
	app.Static("/public", "./public", fiber.Static{
		Compress:  true,
		ByteRange: true,
	})

	api := app.Group("/api")

	api.Get("/web-push-public-key", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"publicKey": webPushPublicKey(),
		})
	})

	USER = api.Group("/user")
	SetupUserRoutes()

	SHOP = api.Group("/shop")
	SetupShopRoutes()

	NOTIFICATION = api.Group("/notification")
	SetupNotificationRoutes()

	TRACKING = api.Group("/tracking")
	SetupTrackingRoutes()
}
