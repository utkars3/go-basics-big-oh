package main

import (
	"fmt"
	"log"
	"microservices/user-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user/:id", handler.GetUserHandler)

	fmt.Println("User Service is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting User Service: ", err)
	}
}
