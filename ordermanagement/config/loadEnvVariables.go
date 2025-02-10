package config

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	envFilePath := filepath.Join("config", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
