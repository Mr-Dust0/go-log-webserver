package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"webserver/types"
)

func (db *Database) InsertTestData() {
	println("Inserting test data")

	// Hash Sample user password to store in database
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password")
	}

	sampleUser := types.User{
		Email:    "test@test.com",
		Username: "test",
		Password: string(passwordHash),
	}

	err = db.Instance.Create(&sampleUser).Error
	if err != nil {
		fmt.Println("Already created sample user")
		return
	}

	// Create sample log entries to show how logs are displayed in the application
	logEntries := []types.LogEntry{
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
	if err := db.Instance.Create(&logEntries).Error; err != nil {
		log.Fatal("Error inserting records for log entries")
	}
}
