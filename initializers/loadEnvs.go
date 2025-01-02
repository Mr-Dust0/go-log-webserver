package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

var EnvFile map[string]string

func LoadEnvs() {
	var err error
	EnvFile, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
