package repositories

import (
	"ordermanagement/config"
	"ordermanagement/internal/models"
)

func CreateUser(user *models.User) error {
	result := config.DBPostgreSQL.Create(&user)
	return result.Error
}

// UserRepository defines the methods that the service layer can call.
type UserRepository interface {
	GetUserByID(userID string) (*models.User, error)
}

// UserRepositoryImpl is the concrete implementation of the UserRepository interface.
type UserRepositoryImpl struct{}

// NewUserRepository creates a new instance of UserRepositoryImpl.
func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// GetUserByID retrieves a user by ID from the database.
func (repo *UserRepositoryImpl) GetUserByID(userID string) (*models.User, error) {
	// Actual implementation with ORM (e.g., GORM)
	return nil, nil
}
