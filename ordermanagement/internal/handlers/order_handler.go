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

	userID, exists := c.Get("userID")
	if !exists {
		// Handle case where userID is not present in the context
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	order.UserID = userID.(uuid.UUID)

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderWithUser, err := services.CreateOrderService(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderWithUser)

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
