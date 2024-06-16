package router

import (
	models "go-authentication-boilerplate/models"
	util "go-authentication-boilerplate/util"
	auth "go-authentication-boilerplate/auth"

	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupPostRoutes() {
	// set up
	privPost := POST.Group("/private")
	privPost.Use(auth.SecureAuth()) // middleware to secure all routes for this group
	privPost.Post("/set", HandleSetPost)
	privPost.Get("/all", HandleGetPosts)
	privPost.Get("/:id", HandleGetPost)
}

func HandleGetPosts(c *fiber.Ctx) error {
	posts, err := util.GetPosts(c.Locals("id").(string))
	if err != nil {
		log.Printf("[ERROR] Error in get posts API: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error getting posts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Posts fetched successfully",
		"posts": posts,
	})
}


func HandleGetPost(c *fiber.Ctx) error {
	postID := c.Params("id")
	post, err := util.GetPost(postID)
	if err != nil {
		log.Printf("[ERROR] Error in get post API: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Post with this ID doesn't exist",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message": "Post fetched successfully",
		"post": post,
	})
}

func HandleSetPost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"message":   "Couldn't parse the input",
		})
	}

	if post.ID == "" {
		post.OwnerID = c.Locals("id").(string)
	} else {
		// get the post from the database
		dbPost, err := util.GetPost(post.ID)
		if err != nil {
			log.Printf("[ERROR] Error in get post API: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"message":   "Post with this ID doesn't exist",
			})
		}

		// check if the user is the owner of the post
		if dbPost.OwnerID != c.Locals("id").(string) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"message":   "You are not the owner of this post",
			})
		}
	}

	savedPost, err := util.SetPost(post)
	if err != nil {
		log.Printf("[ERROR] Error in set post API: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"message":   "Error setting post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message":   "Post set successfully",
		"post": savedPost,
	})
}

