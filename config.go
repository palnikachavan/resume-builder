package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

// GetEnv fetches the environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}
