package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBPostgreSQL *gorm.DB

func ConnectToPostgreSQL() error {
	var err error
	dsn := os.Getenv("DB_URL")
	DBPostgreSQL, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB")
		return err
	}
	// fmt.Println("hello mongo ")
	return nil
}
