package repositories

import (
	"context"
	"fmt"
	"ordermanagement/config"
	"ordermanagement/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// OrderRepository handles order database operations
type OrderRepository struct {
	DB *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		DB: config.DBPostgreSQL, // Assuming DBPostgreSQL is initialized
	}
}

// GetOrderByID fetches an order by its ID from PostgreSQL
func (repo *OrderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := repo.DB.Preload("User").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func FetchOrderHistory(orderID string) (map[string]interface{}, error) {
	collection := config.MongoDB.Collection("order_history")
	filter := bson.M{"order_id": orderID}

	var orderHistory map[string]interface{}
	err := collection.FindOne(
		context.TODO(),
		filter,
	).Decode(&orderHistory)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return orderHistory, nil
}

func FetchOrdersWithPagination(userID string, page, limit int) ([]models.Order, int64, error) {
	// Start query with filtering applied first
	query := config.DBPostgreSQL.Model(&models.Order{}).Order("created_at DESC")

	// Apply user_id filter before pagination
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Get total count of filtered records for frontend pagination
	var totalCount int64
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	var orders []models.Order
	err = query.Limit(limit).Offset((page - 1) * limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func GetOrderByID(db *gorm.DB, orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetInventoryByProductID retrieves inventory for a specific product.
func GetInventoryByProductID(db *gorm.DB, productID uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	if err := db.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

// UpdateInventoryStock updates the inventory for a product.
func UpdateInventoryStock(db *gorm.DB, inventory *models.Inventory) error {
	return db.Save(inventory).Error
}

// SaveOrder updates an order's status.
func SaveOrder(db *gorm.DB, order *models.Order) error {
	return db.Save(order).Error
}

func CreateOrder(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

