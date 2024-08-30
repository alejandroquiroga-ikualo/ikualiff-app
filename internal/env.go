package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	API_KEY             EnvKey = "API_KEY"
	API_JWT_KEY                = "API_JWT_KEY"
	VERIFF_URL                 = "VERIFF_URL"
	VERIFF_IDV_API_KEY         = "VERIFF_IDV_API_KEY"
	VERIFF_POA_API_KEY         = "VERIFF_POA_API_KEY"
	VERIFF_CALLBACK_URL        = "VERIFF_CALLBACK_URL"
	DATABASE_URL               = "DATABASE_URL"
)

func GetEnv() map[EnvKey]string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file!")
	}

	envRegistry := []EnvKey{
		API_KEY, 
		API_JWT_KEY, 
		VERIFF_URL, 
		VERIFF_IDV_API_KEY, 
		VERIFF_POA_API_KEY,
		VERIFF_CALLBACK_URL, 
		DATABASE_URL,
	}
	result := make(map[EnvKey]string)

	for _, registry := range envRegistry {
		envValue := os.Getenv(string(registry))
		if envValue == "" {
			log.Fatalf("%q env variable not set!", string(registry))
		}

		result[registry] = envValue
	}

	return result
}
