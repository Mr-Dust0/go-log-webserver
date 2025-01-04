package initializers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"webserver/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database connection
func InitDatabase() {
	var err error
	// Get the path of executable that is being executed
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// Get the directory where the executable is
	exPath := filepath.Dir(ex)
	// Open database stored in current directory
	DB, err = gorm.Open(sqlite.Open(exPath+"/tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	// Create table for LogEntry is there is none
	err = DB.AutoMigrate(&models.LogEntry{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
	// Create table for User if there is none
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
}
func InsertTestData() {
	// Hash Sample user password to store in database
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	sampleUser := models.User{
		Email:    "test@test.com",
		Username: "test",
		Password: string(passwordHash)}
	if err := DB.Create(&sampleUser).Error; err != nil {
		fmt.Println("Already created sample user")
		return

	}
	// Create sample log entries to show how logs are displayed in the application
	logEntries := []models.LogEntry{
		{
			TimeStamp:       time.Date(2025, 1, 2, 15, 5, 47, 42288622, time.UTC),
			TimeStampClosed: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), // Invalid date
			HostName:        "debian",
			FileName:        "/mnt/nfs/testdoc.pdf",
		},
		{
			TimeStamp:       time.Date(2025, 1, 3, 9, 50, 29, 657848343, time.UTC),
			TimeStampClosed: time.Date(2025, 1, 3, 9, 51, 47, 715981132, time.UTC),
			HostName:        "raspberrypi",
			FileName:        "/home/pi/009085-A Drawing v1.pdf",
		},
		{
			TimeStamp:       time.Date(2025, 1, 3, 9, 57, 23, 461050454, time.UTC),
			TimeStampClosed: time.Date(2025, 1, 3, 10, 2, 54, 111097126, time.UTC),
			HostName:        "raspberrypi",
			FileName:        "/home/pi/testdoc.pdf",
		},
		{
			TimeStamp:       time.Date(2025, 1, 3, 10, 2, 57, 659195934, time.UTC),
			TimeStampClosed: time.Date(2025, 1, 3, 10, 23, 22, 617727232, time.UTC),
			HostName:        "raspberrypi",
			FileName:        "/home/pi/009085-A Drawing v1.pdf",
		},
		{
			TimeStamp:       time.Date(2025, 1, 3, 10, 23, 33, 658710661, time.UTC),
			TimeStampClosed: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), // Invalid date
			HostName:        "raspberrypi",
			FileName:        "/home/pi/testdoc.pdf",
		},
	}

	// Insert the records into the database
	if err := DB.Create(&logEntries).Error; err != nil {
		log.Fatal("Error inserting records for log enteries")
	}
}
