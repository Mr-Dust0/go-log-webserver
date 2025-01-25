package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"webserver/types"
)

var SCHEMES_TO_MIGRATE = []interface{}{
	&types.LogEntry{},
	&types.User{},
}

// DB is a struct that represents a high-level database package
type Database struct {
	credentials Credentials
	Instance    *gorm.DB
}

// Credentials stores the database configuration
type Credentials struct {
	URI string
}

// Start is a function tha fetches credentials and initiates the connection to the database
func NewDatabaseFromEnv() *Database {
	log.Println("Connecting to Database...")

	// Get the database credentials from the environment
	DB_NAME, _ := os.LookupEnv("DB_NAME")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	db := &Database{
		credentials: Credentials{
			URI: filepath.Dir(ex) + DB_NAME,
		},
	}

	db.connect()

	err = db.migrate()
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database ready for use!")
	log.Println("Database location: ", db.credentials.URI)

	return db
}

func (db *Database) connect() {
	sqlDB, err := gorm.Open(sqlite.Open(db.credentials.URI), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("Failed to connect to database: %v", err.Error()))
	}

	db.Instance = sqlDB
}

// Initialize the database connection
func (db *Database) migrate() error {
	for _, scheme := range SCHEMES_TO_MIGRATE {
		err := db.Instance.AutoMigrate(scheme)
		if err != nil {
			log.Println("Failed to migrate scheme: ", err)
			return err
		}
	}

	return nil
}
