package main

import (
	"fmt"
	"log"
	"microservices/notification-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/notify", handler.NotifyHandler)

	fmt.Println("Notification Service is running on port 8081...")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Error starting Notification Service: ", err)
	}
}
