package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ordermanagement/config"
	"ordermanagement/internal/models"
	"ordermanagement/internal/repositories"
	"ordermanagement/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// PlaceOrder places a new order with inventory check and MongoDB logging
func PlaceOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	// Check inventory using ProductID
	var inventory models.Inventory
	if err := tx.Where("id = ?", order.ProductID).First(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Ensure stock is available
	if inventory.Stock <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Out of stock"})
		return
	}

	// Deduct stock & create order
	inventory.Stock--
	if err := tx.Save(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		return
	}

	// Create the order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	var orderWithUser models.Order
	if err := tx.Preload("User").First(&orderWithUser, order.ID).Error; err != nil {
		log.Println("Error fetching order with user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order with user details"})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Log order in MongoDB
	logOrderHistory(order.ID, "pending")

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": orderWithUser})
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

// GetOrder retrieves an order by its ID from PostgreSQL
func GetOrder(c *gin.Context) {
	// Extract order ID from the URL parameter
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Convert string orderID to uuid
	parsedOrderID, err := uuid.Parse(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Fetch order details from PostgreSQL
	orderRepo := repositories.NewOrderRepository()
	order, err := orderRepo.GetOrderByID(parsedOrderID)
	if err != nil {
		log.Println("Error fetching order from PostgreSQL:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Return the order details
	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// func CancelOrderHandler(c *gin.Context) {
// 	orderID := c.Param("id")
// 	parsedOrderID, err := uuid.Parse(orderID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
// 		return
// 	}
// 	order, err := services.CancelOrderService(config.DBPostgreSQL, parsedOrderID)
// 	if err != nil {
// 		if fmt.Sprintf("%T", err) == "*gorm.ErrRecordNotFound" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with ID %s not found", orderID)})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to cancel order: %v", err)})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Order cancelled successfully",
// 		"order":   order,
// 	})
// }

// CancelOrder cancels an existing order by updating its status to 'cancelled' if it's in 'pending' state
func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Retrieve the order by its ID
	var order models.Order
	if err := config.DBPostgreSQL.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with ID %s not found", orderID)})
		return
	}

	// Check if the order is in the 'pending' state
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Order cannot be cancelled because it is in '%s' state", order.Status),
		})
		return
	}

	// 	//--------------------add back to inventory--------------------------
	var inventory models.Inventory
	if err := config.DBPostgreSQL.Where("id = ?", order.ProductID).First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product with ID %s not found in inventory", order.ProductID)})
		return
	}
	inventory.Stock++
	if err := config.DBPostgreSQL.Save(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update inventory for product %s", order.ProductID)})
		return
	}

	// Log the status before cancelling the order
	logOrderStatusChange(order.ID, order.Status, "cancelled")

	// Update the status of the order to 'cancelled'
	order.Status = "cancelled"
	if err := config.DBPostgreSQL.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel the order", "details": err.Error()})
		return
	}

	// Return success response with the updated order details
	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"order":   order,
	})
}

// logOrderStatusChange logs the status change in MongoDB
func logOrderStatusChange(orderID uuid.UUID, oldStatus string, newStatus string) {
	timestamp := time.Now()

	_, err := config.MongoDB.Collection("order_history").UpdateOne(
		context.TODO(),
		bson.M{"order_id": orderID.String()},
		bson.M{
			"$set": bson.M{
				"status": "cancelled",
			},
			"$push": bson.M{
				"log": bson.M{
					"status":    "cancelled",
					"timestamp": timestamp,
				},
			},
		},
	)
	if err != nil {
		log.Println("Failed to update MongoDB:", err)
	}
}

func GetOrdersWithPagination(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID := c.Query("user_id") // Get user_id from query params

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Call service to fetch orders with pagination
	ordersResponse, err := services.GetOrdersWithPaginationService(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ordersResponse)
}

// completed
func HistoryOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	orderHistory, err := services.GetOrderHistory(orderID)
	if err != nil {
		if err.Error() == "Order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching order history"})
			log.Println("Error fetching order history:", err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_history": orderHistory})

}
