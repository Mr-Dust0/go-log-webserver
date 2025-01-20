package utils

import (
	"sync"

	"github.com/joho/godotenv"
)

var (
	envVarsOnce sync.Once
	envVars     map[string]string
)

func GetDotEnvVariables() map[string]string {
	envVarsOnce.Do(func() {
		var err error
		envVars, err = godotenv.Read()
		if err != nil {
			panic(err)
		}
	})

	return envVars
}

func IsTestEnv() bool {
	env := GetDotEnvVariables()
	return env["TEST_ENV"] == "true"
}
