package main

import (
	"fmt"
	"log"
	"ordermanagement/config"
	"ordermanagement/middleware"
	"ordermanagement/routes"
	"ordermanagement/workers"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToPostgreSQL()
	config.ConnectToMongoDB()
}

func main() {
	// Initialize the Gin router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	// Setup the routes
	routes.SetupRoutes(r)

	go workers.ProcessOrders(config.DBPostgreSQL)

	// Start the Gin HTTP server
	fmt.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
