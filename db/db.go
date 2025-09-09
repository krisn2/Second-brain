package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dburl := ``
	db, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	DB = db
}
