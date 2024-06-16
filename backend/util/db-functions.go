package util

import (
	"log"
	db "go-authentication-boilerplate/database"
	models "go-authentication-boilerplate/models"
)

func GetPosts(ownerID string) ([]models.Post, error) {
	posts := []models.Post{}
	txn := db.DB.Where("owner_id = ?", ownerID).Find(&posts)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting posts: %v", txn.Error)
		return posts, txn.Error
	}
	return posts, nil
}


func GetPost(id string) (*models.Post, error) {
	post := new(models.Post)
	txn := db.DB.Where("id = ?", id).First(post)
	if txn.Error != nil {
		log.Printf("[ERROR] Error getting post: %v", txn.Error)
		return post, txn.Error
	}
	return post, nil
}

func SetPost(post *models.Post) (*models.Post, error) {
	// check if post with ID exists
	if post.ID == "" {
		post.CreatedAt = db.DB.NowFunc().String()
		post.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Create(post)
		if txn.Error != nil {
			log.Printf("[ERROR] Error creating post: %v", txn.Error)
			return post, txn.Error
		}
	} else {
		post.UpdatedAt = db.DB.NowFunc().String()
		txn := db.DB.Save(post)
		if txn.Error != nil {
			log.Printf("[ERROR] Error saving post: %v", txn.Error)
			return post, txn.Error
		}
	}

	return post, nil
}
