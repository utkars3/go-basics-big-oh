package repository

import (
	"errors"
	model "microservices/user-service/internal/models"
)

// Simulated function to get a user from the database
func GetUserByID(userID int) (model.User, error) {
	// In a real scenario, replace this with a database query
	if userID == 1 {
		return model.User{ID: 1, Name: "John Doe"}, nil
	} else if userID == 2 {
		return model.User{ID: 2, Name: "Utkarsh Kesharwani"}, nil
	}
	return model.User{}, errors.New("user not found")
}
