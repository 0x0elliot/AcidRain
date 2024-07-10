package router

import (
	util "go-authentication-boilerplate/util"
	auth "go-authentication-boilerplate/auth"

	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupShopRoutes() {
	// set up
	privShop := SHOP.Group("/private")
	privShop.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privShop.Get("/all", HandleGetAllAccessibleShops)
}

func HandleGetAllAccessibleShops(c *fiber.Ctx) error {
	shops, err := util.GetShops(c.Locals("id").(string))
	if err != nil {
		log.Printf("[ERROR] Error in get shops API: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting shops",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Shops fetched successfully",
		"shops": shops,
	})
}


