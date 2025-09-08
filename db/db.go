package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dburl := `postgresql://neondb_owner:npg_5TlhJLw2UHAV@ep-icy-dust-adj47zfy-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require`
	db, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	DB = db
}
