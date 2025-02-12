package repository

import (
	"log"
	"microservices/notification-service/internal/models"
)

func SaveNotification(notification models.Notification) error {
	log.Printf("Storing notification for UserID %d: %s\n", notification.UserID, notification.Message)
	return nil
}








