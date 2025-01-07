package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetFromEnv gets value of envs
func GetFromEnv(s string) string {
	v := os.Getenv(s)
	if v == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
		v = os.Getenv(s)
	}

	return v
}
