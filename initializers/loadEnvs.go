package initializers

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Enable EnvFile to be accessed from anywhere
var EnvFile map[string]string

func LoadEnvs() {
	var err error
	// load .env file that config in
	EnvFile, err = godotenv.Read()
	if err != nil {
		//log.Fatal("Error loading .env file")
		fmt.Println("Could not find env file")

	}

}
