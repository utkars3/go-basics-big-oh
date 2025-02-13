package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	model "microservices/user-service/internal/models"
	"microservices/user-service/internal/repository"
	"net/http"
)

// Fetch user from repository
func GetUserByID(userID int) (model.User, error) {
	return repository.GetUserByID(userID)
}

// Send notification to the user
func SendNotification(user model.User) {
	notificationURL := "http://localhost:8081/notification"

	notificationData := map[string]interface{}{
		"user_id": user.ID,
		"message": fmt.Sprintf("Hello, %s! You have a new message.", user.Name),
	}

	data, err := json.Marshal(notificationData)
	if err != nil {
		log.Println("Error marshalling notification data:", err)
		return
	}

	resp, err := http.Post(notificationURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error sending notification request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Notification Response:", string(body))
}
