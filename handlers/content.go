package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krisn2/second-brain/db"
	"github.com/krisn2/second-brain/models"
)

func AddContent(c *gin.Context) {

	var body struct {
		Title string   `json:"title"`
		Type  string   `json:"type"`
		Link  string   `json:"link"`
		Tags  []string `json:"tags"`
	}

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid Input"})
		return
	}

	userID, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uuidUser, _ := uuid.Parse(userID.(string))

	var tags []models.Tag
	for _, t := range body.Tags {
		tag := models.Tag{Title: t}
		db.DB.FirstOrCreate(&tag, models.Tag{Title: t})
		tags = append(tags, tag)
	}

	content := models.Content{
		Link:   body.Link,
		Type:   body.Type,
		Title:  body.Title,
		UserId: uuidUser,
		Tags:   tags,
	}

	if err := db.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": "Brain created successfully",
		},
	)
}

func GetContent(c *gin.Context) {
	var contents []models.Content
	if err := db.DB.Preload("Tags").Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
	})
}
