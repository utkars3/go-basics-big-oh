package services

import (
	"ordermanagement/internal/models"
	"ordermanagement/internal/repositories"
)

func CreateInventoryService(item *models.Inventory) (*models.Inventory, error) {
	// Validate the inventory item
	if err := item.ValidateInventory(); err != nil {
		return nil, err // Return the validation error
	}

	// Call the repository to create or update the inventory item
	return repositories.CreateInventory(item)
}

func GetAllInventory() ([]models.Inventory, error) {
	return repositories.FetchAllInventory()
}
