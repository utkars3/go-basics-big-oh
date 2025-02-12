package service

import (
	"log"
	"microservices/notification-service/internal/models"
	"microservices/notification-service/internal/repository"
)


func SendNotification(notification models.Notification) error {
	
	log.Printf("Sending notification to UserID %d: %s\n", notification.UserID, notification.Message)

	err := repository.SaveNotification(notification)
	if err != nil {
		return err
	}			
				
	return nil
}
