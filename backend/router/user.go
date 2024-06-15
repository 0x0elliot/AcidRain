package router

import (
	db "go-authentication-boilerplate/database"
	"go-authentication-boilerplate/models"
	auth "go-authentication-boilerplate/auth"
	"log"
	"time"
	// "golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(db.PRIVKEY)

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	// USER.Post("/signup", CreateUser)              // Sign Up a user
	// USER.Post("/signin", LoginUser)               // Sign In a user
	USER.Get("/get-access-token", GetAccessToken) // returns a new access_token
	USER.Post("/passwordless-login", HandlePasswordLessLogin)

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privUser.Get("/getinfo", GetUserData)

}

func HandlePasswordLessLogin(c *fiber.Ctx) error {
	type LoginInput struct {
		Email string `json:"email"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "message": "Please review your input"})
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

	_, accessToken := auth.GenerateAccessClaims(refreshClaims.Issuer)

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

	return c.JSON(u)
}
