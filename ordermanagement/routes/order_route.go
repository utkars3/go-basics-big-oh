package routes

import (
	"ordermanagement/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes initializes order-related routes
func SetupOrderRoutes(api *gin.RouterGroup) {
	orderRoutes := api.Group("/orders")
	{
		orderRoutes.POST("/", handlers.PlaceOrder)             // Place a new order
		orderRoutes.GET("/:id", handlers.GetOrder)             // Fetch order details
		orderRoutes.GET("/", handlers.GetOrdersWithPagination) // Fetch order details
		orderRoutes.PATCH("/:id/cancel", handlers.CancelOrder) // Cancel an order
		orderRoutes.GET("/:id/history", handlers.HistoryOrder) // Cancel an order
	}
}
