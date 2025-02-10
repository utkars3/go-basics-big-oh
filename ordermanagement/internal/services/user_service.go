package services

import (
	"errors"
	"ordermanagement/internal/models"
	"ordermanagement/internal/repositories"
)

func CreateUserService(user *models.User) error {
	if err := user.ValidateUser(); err != nil {
		return errors.New("validation failed")
	}
	return repositories.CreateUser(user)
}
