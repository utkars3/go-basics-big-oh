package handler

import (
	"fmt"
	"microservices/notification-service/internal/models"
	"microservices/notification-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotifyHandler handles incoming notification requests
func NotifyHandler(c *gin.Context) {
	var notification models.Notification
	fmt.Println("-------------i am coming-----------")
	// Parse incoming JSON request
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Process notification
	if err := service.SendNotification(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	// Send success response
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}
