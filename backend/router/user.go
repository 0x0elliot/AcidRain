package router

import (
	db "go-authentication-boilerplate/database"
	"go-authentication-boilerplate/models"
	"go-authentication-boilerplate/util"
	// "math/rand"
	"time"
	"log"
	// "golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(db.PRIVKEY)

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	// USER.Post("/signup", CreateUser)              // Sign Up a user
	// USER.Post("/signin", LoginUser)               // Sign In a user
	USER.Post("/whatsapp-login", HandleWhatsAppLogin) // Sign In a user using WhatsApp
	USER.Get("/get-access-token", GetAccessToken) // returns a new access_token

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(util.SecureAuth()) // middleware to secure all routes for this group
	privUser.Get("/user", GetUserData)

}

// CreateUser route registers a User into the database

// func CreateUser(c *fiber.Ctx) error {
// 	u := new(models.User)

// 	if err := c.BodyParser(u); err != nil {
// 		return c.JSON(fiber.Map{
// 			"error": true,
// 			"input": "Please review your input",
// 		})
// 	}

// 	// validate if the email, username and password are in correct format
// 	errors := util.ValidateRegister(u)
// 	if errors.Err {
// 		return c.JSON(errors)
// 	}

// 	if count := db.DB.Where(&models.User{Email: u.Email}).First(new(models.User)).RowsAffected; count > 0 {
// 		errors.Err, errors.Email = true, "Email is already registered"
// 	}
// 	if count := db.DB.Where(&models.User{Username: u.Username}).First(new(models.User)).RowsAffected; count > 0 {
// 		errors.Err, errors.Username = true, "Username is already registered"
// 	}
// 	if errors.Err {
// 		return c.JSON(errors)
// 	}

// 	// Hashing the password with a random salt
// 	password := []byte(u.Password)
// 	hashedPassword, err := bcrypt.GenerateFromPassword(
// 		password,
// 		rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost,
// 	)

// 	if err != nil {
// 		panic(err)
// 	}
// 	u.Password = string(hashedPassword)

// 	if err := db.DB.Create(&u).Error; err != nil {
// 		return c.JSON(fiber.Map{
// 			"error":   true,
// 			"general": "Something went wrong, please try again later. 😕",
// 		})
// 	}

// 	// setting up the authorization cookies
// 	accessToken, refreshToken := util.GenerateTokens(u.UUID.String())
// 	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
// 	c.Cookie(accessCookie)
// 	c.Cookie(refreshCookie)

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 	})
// }

// LoginUser route logins a user in the app
// func LoginUser(c *fiber.Ctx) error {
// 	type LoginInput struct {
// 		Identity string `json:"identity"`
// 		Password string `json:"password"`
// 	}

// 	input := new(LoginInput)

// 	if err := c.BodyParser(input); err != nil {
// 		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
// 	}

// 	u := new(models.User)
// 	if res := db.DB.Where(
// 		&models.User{Email: input.Identity}).Or(
// 		&models.User{Username: input.Identity},
// 	).First(&u); res.RowsAffected <= 0 {
// 		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
// 	}

// 	// Comparing the password with the hash
// 	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
// 		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
// 	}

// 	// setting up the authorization cookies
// 	accessToken, refreshToken := util.GenerateTokens(u.UUID.String())
// 	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
// 	c.Cookie(accessCookie)
// 	c.Cookie(refreshCookie)

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 	})
// }

func HandleWhatsAppLogin(c *fiber.Ctx) error {
	type LoginInput struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		log.Printf("[ERROR] Failed to parse the input: %v", err)
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	if len(input.OTP) == 0 {
		u := new(models.User)
		if res := db.DB.Where(
			&models.User{Phone: input.Phone},
		).First(&u); res.RowsAffected <= 0 {
			if err := db.DB.Create(&models.User{Phone: input.Phone}).Error; err != nil {
				log.Printf("[ERROR] Failed to create user: %v", err)
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{"error": true, "general": "Something went wrong, please try again later!"})
			}
		}

		if !util.IsValidPhone(input.Phone) {
			log.Printf("[ERROR] Invalid phone number: %v", input.Phone)
			return c.JSON(fiber.Map{"error": true, "phone": "Invalid phone number"})
		}

		// send OTP to the phone number
		err := util.GenerateOTP(input.Phone)
		if err != nil {
			log.Printf("[ERROR] Failed to send OTP: %v", err)
			return c.JSON(fiber.Map{"error": true, "general": "Failed to send OTP"})
		}
		return c.JSON(fiber.Map{"error": false, "general": "OTP sent successfully"})
	}

	// verify the OTP
	return c.JSON(fiber.Map{"error": false, "general": "OTP verified successfully"})
}

// GetAccessToken generates and sends a new access token iff there is a valid refresh token
func GetAccessToken(c *fiber.Ctx) error {
	type RefreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}

	reToken := new(RefreshToken)
	if err := c.BodyParser(reToken); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
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

	_, accessToken := util.GenerateAccessClaims(refreshClaims.Issuer)

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
	if res := db.DB.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}
