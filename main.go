package main

import (
	"github.com/gin-gonic/gin"
	"github.com/krisn2/second-brain/db"
	"github.com/krisn2/second-brain/handlers"
	"github.com/krisn2/second-brain/middleware"
	"github.com/krisn2/second-brain/models"
)

func main() {
	router := gin.Default()
	db.Connect()
	if err := db.DB.AutoMigrate(&models.User{}, &models.Content{}, &models.Tag{}); err != nil {
		panic(err)
	}

	router.POST("/api/v1/signup", handlers.Register)
	router.POST("/api/v1/login", handlers.Login)

	content := router.Group("/api/v1/content", middleware.AuthMiddleware())
	{
		content.GET("", handlers.GetContent)
		content.POST("", handlers.AddContent)
	}

	router.Run()
}
