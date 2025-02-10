package main

import (
	"log"
	"ordermanagement/config"
	"ordermanagement/internal/models"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToPostgreSQL()
}

func main() {
	if err := config.DBPostgreSQL.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration failed for User model: %v", err)
	}
	log.Println("User model migrated successfully.")

	if err := config.DBPostgreSQL.AutoMigrate(&models.Inventory{}); err != nil {
		log.Fatalf("Migration failed for Inventory model: %v", err)
	}
	log.Println("Inventory model migrated successfully.")

	if err := config.DBPostgreSQL.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("Migration failed for Order model: %v", err)
	}
	log.Println("Order model migrated successfully.")
}
