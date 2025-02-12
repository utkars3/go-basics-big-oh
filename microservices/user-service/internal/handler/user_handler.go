package handler

import (
	"log"
	"microservices/user-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserHandler handles fetching a user and sending a notification
func GetUserHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// fmt.Println(user)
	c.JSON(http.StatusOK, user)

	// Send notification after responding to the client
	go service.SendNotification(user) // Run asynchronously
	log.Println("Notification triggered for user:", user.Name)
}
