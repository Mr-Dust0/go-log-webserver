package initializers

import (
	"log"
	"webserver/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database connection
func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("/home/john/Downloads/dad-project/test/build/pdf_tracker/tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	err = DB.AutoMigrate(&models.LogEntry{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
}
