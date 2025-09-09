package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/krisn2/second-brain/db"
	"github.com/krisn2/second-brain/handlers"
	"github.com/krisn2/second-brain/middleware"
	"github.com/krisn2/second-brain/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		log.Fatal(err)
	}
	router := gin.Default()
	db.Connect(os.Getenv("DATABASE_URL"))
	if err := db.DB.AutoMigrate(&models.User{}, &models.Content{}, &models.Tag{}); err != nil {
		panic(err)
	}

	router.POST("/api/v1/signup", handlers.Register)
	router.POST("/api/v1/login", handlers.Login)

	content := router.Group("/api/v1/content", middleware.AuthMiddleware())
	{
		content.GET("", handlers.SearchBrain)
		content.POST("", handlers.AddContent)
		content.DELETE("", handlers.DeleteContent)
		content.GET("", handlers.GetContent)
	}

	router.Run()
}
