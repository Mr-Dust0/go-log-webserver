package initializers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"webserver/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database connection
func InitDatabase() {
	var err error
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	DB, err = gorm.Open(sqlite.Open(exPath+"/tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	err = DB.AutoMigrate(&models.LogEntry{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
}
func InsertTestData() {
	sampleUser := models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: "password"}
	if err := DB.Create(&sampleUser).Error; err != nil {
		log.Fatal("Could not create the sample user")
		return

	}
}
