package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ordermanagement/config"
	"ordermanagement/internal/models"
	"ordermanagement/internal/repositories"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type OrdersResponse struct {
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int64          `json:"total"`
	TotalPages int64          `json:"totalPages"`
	Orders     []models.Order `json:"orders"`
}

func GetOrderHistory(orderID string) (map[string]interface{}, error) {
	// Call the repository to fetch order history by orderID
	orderHistory, err := repositories.FetchOrderHistory(orderID)
	if err != nil {
		log.Println("Error in service:", err)
		return nil, err
	}
	return orderHistory, nil
}

func GetOrdersWithPaginationService(userID string, page, limit int) (*OrdersResponse, error) {
	// Call repository to fetch orders with pagination
	orders, totalCount, err := repositories.FetchOrdersWithPagination(userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}

	// Calculate total pages
	totalPages := (totalCount + int64(limit) - 1) / int64(limit)

	// Return paginated response
	return &OrdersResponse{
		Page:       page,
		Limit:      limit,
		Total:      totalCount,
		TotalPages: totalPages,
		Orders:     orders,
	}, nil
}

func CancelOrderService(db *gorm.DB, orderID uuid.UUID) (*models.Order, error) {
	order, err := repositories.GetOrderByID(db, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.Status != "pending" {
		return nil, fmt.Errorf("order cannot be cancelled because it is in '%s' state", order.Status)
	}

	inventory, err := repositories.GetInventoryByProductID(db, order.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found in inventory: %w", err)
	}

	inventory.Stock++
	if err := repositories.UpdateInventoryStock(db, inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	order.Status = "cancelled"
	if err := repositories.SaveOrder(db, order); err != nil {
		return nil, fmt.Errorf("failed to cancel the order: %w", err)
	}

	return order, nil
}

func CreateOrderService(order *models.Order) (*models.Order, error) {
	// Begin transaction
	tx := config.DBPostgreSQL.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the user exists
	var user models.User
	if err := tx.Where("id = ?", order.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user not found")
	}

	// Check inventory using ProductID
	var inventory models.Inventory
	if err := tx.Where("id = ?", order.ProductID).First(&inventory).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("product not found")
	}

	// Ensure stock is available
	if inventory.Stock <= 0 {
		tx.Rollback()
		return nil, errors.New("out of stock")
	}

	// Deduct stock & create order
	inventory.Stock--
	if err := repositories.SaveInventory(tx, &inventory); err != nil {
		tx.Rollback()
		return order, errors.New("failed to update inventory")
	}

	// Create the order using the repository
	if err := repositories.CreateOrder(tx, order); err != nil {
		tx.Rollback()
		return order, errors.New("failed to place order")
	}

	var orderWithUser models.Order
	fmt.Println("------------------------------------", order.ID)
	if err := tx.Preload("User").First(&orderWithUser, order.ID).Error; err != nil {
		log.Println("Error fetching order with user:", err)
		return order, errors.New("failed to fetch order with user details")
	}

	tx.Commit()

	logOrderHistory(order.ID, "pending")

	return &orderWithUser, nil
}

// logOrderHistory logs order history in MongoDB
func logOrderHistory(orderID uuid.UUID, status string) {
	// Get the current timestamp
	timestamp := time.Now()
	orderIDStr := orderID.String()
	// Create the order history log entry
	orderHistory := bson.M{
		"order_id": orderIDStr,
		"status":   status,
		"log": []bson.M{
			{
				"status":    status,
				"timestamp": timestamp, // Add the timestamp
			},
		},
	}

	// Insert the order history into MongoDB
	_, err := config.MongoDB.Collection("order_history").InsertOne(context.Background(), orderHistory)
	if err != nil {
		log.Println("Error logging order history:", err)
	}
}
