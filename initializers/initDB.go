package initializers

import (
	"log"
	"os"
	"personal-blog/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	// Initialize database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.Post{})

	// Set global variable
	DB = db
}
