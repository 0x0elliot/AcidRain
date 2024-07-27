package router

import (
	db "go-authentication-boilerplate/database"
	"go-authentication-boilerplate/models"
	auth "go-authentication-boilerplate/auth"
	util "go-authentication-boilerplate/util"
	"log"
	"time"
	// "golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(db.PRIVKEY)

func SetupUserRoutes() {
	USER.Get("/get-access-token", GetAccessToken) // returns a new access_token
	USER.Post("/passwordless-login", HandlePasswordLessLogin)
	USER.Get("/shopify-oauth", HandleRedirectToShopifyOAuth)

	USER.Post("/logout", HandleLogout)
	USER.Get("/logout", HandleLogout)
	USER.Post("/verify-token", HandleVerifyToken)

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privUser.Get("/getinfo", GetUserData)
	privUser.Get("/shopify/callback", HandleShopifyOauthCallback)
}

func HandleVerifyToken(c *fiber.Ctx) error {
	// check if the tokens sent are valid
	type Tokens struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	tokens := new(Tokens)
	if err := c.BodyParser(tokens); err != nil {
		log.Printf("[ERROR] Couldn't parse the input: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Please review your input"})
	}

	accessToken := tokens.AccessToken
	refreshToken := tokens.RefreshToken

	accessClaims := new(models.Claims)
	_, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
	})

	if err != nil {
		log.Printf("[ERROR] Couldn't parse access token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "message": "Invalid access token"})
	}

	refreshClaims := new(models.Claims)
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
	})

	if err != nil {
		log.Printf("[ERROR] Couldn't parse refresh token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "message": "Invalid refresh token"})
	}

	// set cookies if valid
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: accessToken,
		Expires: time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure: true,
	})

	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		Expires: time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure: true,
	})

	return c.JSON(fiber.Map{"message": "Token verified successfully"})
}


func HandleLogout(c *fiber.Ctx) error {
	c.ClearCookie("access_token", "refresh_token")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}

func HandleShopifyOauthCallback(c *fiber.Ctx) error {
	// get the query params
	code := c.Query("code")
	hmac := c.Query("hmac")
	host := c.Query("host")
	shop := c.Query("shop")
	state := c.Query("state")
	timestamp := c.Query("timestamp")

	if len(code) == 0 || len(hmac) == 0 || len(host) == 0 || len(shop) == 0 || len(state) == 0 || len(timestamp) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid query parameters"})
	}

	ctx := c.Context()

	// get the access token
	accessToken, err := auth.ExchangeToken(ctx, shop, code)
	if err != nil {
		log.Printf("[ERROR] Couldn't get access token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Couldn't get access token"})
	}

	userID := c.Locals("id").(string)

	// return c.JSON(fiber.Map{"access_token": accessToken})

	// check if the shop exists
	s := new(models.Shop)
	if res := db.DB.Where("shop_identifier = ?", shop).First(&s); res.RowsAffected <= 0 {
		// create a new shop
		s.Name = shop
		s.ShopIdentifier = shop
		s.AccessToken = accessToken
		s.OwnerID = userID
		s.Platform = "shopify"
		if err := db.DB.Create(&s).Error; err != nil {
			log.Printf("[ERROR] Couldn't create shop: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Couldn't create shop"})
		}
	} else {
		// update the access token
		s.AccessToken = accessToken
		if err := db.DB.Save(&s).Error; err != nil {
			log.Printf("[ERROR] Couldn't update shop: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Couldn't update shop"})
		}
	}

	// update the user's current shop
	u, err := util.GetUserById(userID)
	if err != nil {
		log.Printf("[ERROR] Couldn't get user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Couldn't get user"})
	}

	u.CurrentShop = *s
	// update the user
	u, err = util.SetUser(u)
	if err != nil {
		log.Printf("[ERROR] Couldn't set user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Couldn't set user"})
	}

	return c.JSON(fiber.Map{"message": "Shopify OAuth callback successful"})
}

func HandleRedirectToShopifyOAuth(c *fiber.Ctx) error {
	shopName := c.Query("shop")

	if len(shopName) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Shop name is required"})
	}

	authURL := auth.GenerateAuthURL(shopName)
	// redirect to the shopify auth url
	return c.JSON(fiber.Map{"auth_url": authURL})
}

func HandlePasswordLessLogin(c *fiber.Ctx) error {
	type LoginInput struct {
		Email string `json:"email"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		log.Printf("[ERROR] Couldn't parse the input: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Please review your input"})
	}

	u := new(models.User)
	if res := db.DB.Where(
		&models.User{Email: input.Email},
	).First(&u); res.RowsAffected <= 0 {
		// create a user with the email
		u.Email = input.Email
		if err := db.DB.Create(&u).Error; err != nil {
			log.Printf("[ERROR] Couldn't create user: %v", err)
			return c.JSON(fiber.Map{"error": true, "message": "Cannot create a user"})
		}
	}

	err := auth.GeneratePasswordLessLink(u)
	if err != nil {
		log.Printf("[INFO] Passwordless login error: %s", err)
		return c.JSON(fiber.Map{"error": true, "message": "Cannot send the passwordless login link"})
	}

	return c.JSON(fiber.Map{"message": "Passwordless login link has been sent to your email"})
}


// GetAccessToken generates and sends a new access token iff there is a valid refresh token
func GetAccessToken(c *fiber.Ctx) error {
	type RefreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}

	reToken := new(RefreshToken)
	if err := c.BodyParser(reToken); err != nil {
		return c.JSON(fiber.Map{"error": true, "message": "Please review your input"})
	}

	refreshToken := reToken.RefreshToken

	refreshClaims := new(models.Claims)
	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if res := db.DB.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
	).First(&models.Claims{}); res.RowsAffected <= 0 {
		// no such refresh token exist in the database
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			// refresh token is expired
			c.ClearCookie("access_token", "refresh_token")
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		// malformed refresh token
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	_, accessToken, err := auth.GenerateAccessClaims(refreshClaims.Issuer)
	if err != nil {
		log.Printf("[ERROR] Couldn't generate access token: %v", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": true, "message": "Couldn't generate access token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}

/*
	PRIVATE ROUTES
*/

// GetUserData returns the details of the user signed in
func GetUserData(c *fiber.Ctx) error {
	id := c.Locals("id")

	u := new(models.User)
	if res := db.DB.Where("id = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "message": "Cannot find the User"})
	}

	// if u.CurrentShop doesn't have a value, attach it to an existing
	// shop owned by the user
	if len(u.CurrentShop.ID) == 0 {
		s := new(models.Shop)
		if res := db.DB.Where("owner_id = ?", u.ID).First(&s); res.RowsAffected > 0 {
			u_ := u
			u_.CurrentShop = *s
			u_.CurrentShopID = &s.ID
			u_, err := util.SetUser(u_)
			if err != nil {
				log.Printf("[ERROR] Couldn't set user: %v", err)
			} else {
				u = u_
			}
		}
	}

	if len(u.CurrentShop.ID) > 0 {
		// get the shop details
		u.CurrentShop.AccessToken = ""
	}

	return c.JSON(u)
}
