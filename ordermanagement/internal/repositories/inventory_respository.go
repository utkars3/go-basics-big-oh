package repositories

import (
	"ordermanagement/config"
	"ordermanagement/internal/models"

	"gorm.io/gorm"
)

func CreateInventory(item *models.Inventory) (*models.Inventory, error) {
	var existingItem models.Inventory

	// Check if the item already exists based on name
	if err := config.DBPostgreSQL.Where("name = ?", item.Name).First(&existingItem).Error; err == nil {
		// If product exists, update the stock count
		existingItem.Stock += item.Stock
		if err := config.DBPostgreSQL.Save(&existingItem).Error; err != nil {
			return nil, err
		}
		// Return the updated item
		return &existingItem, nil
	}

	// If product doesn't exist, create a new entry
	if err := config.DBPostgreSQL.Create(&item).Error; err != nil {
		return nil, err
	}

	// Return the newly created item
	return item, nil
}

func FetchAllInventory() ([]models.Inventory, error) {
	var inventory []models.Inventory
	result := config.DBPostgreSQL.Find(&inventory) // Sorted by latest created items
	return inventory, result.Error
}

type InventoryRepository interface {
	GetInventoryByProductID(productID string) (*models.Inventory, error)
}

// InventoryRepositoryImpl is the concrete implementation of the InventoryRepository interface.
type InventoryRepositoryImpl struct{}

// NewInventoryRepository creates a new instance of InventoryRepositoryImpl.
func NewInventoryRepository() InventoryRepository {
	return &InventoryRepositoryImpl{}
}

// GetInventoryByProductID retrieves inventory by ProductID.
func (repo *InventoryRepositoryImpl) GetInventoryByProductID(productID string) (*models.Inventory, error) {
	// Actual implementation with ORM (e.g., GORM)
	return nil, nil
}

func SaveInventory(tx *gorm.DB, inventory *models.Inventory) error {
	// Save the updated inventory
	return tx.Save(inventory).Error
}
