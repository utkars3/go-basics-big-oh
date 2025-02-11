package routes

import (
	"ordermanagement/internal/handlers"
	"ordermanagement/middleware"
	// "ordermanagement/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes initializes order-related routes
func SetupOrderRoutes(api *gin.RouterGroup) {
	orderRoutes := api.Group("/orders")
	orderRoutes.Use(middleware.JWTAuthMiddleware())
	{
		orderRoutes.POST("/", handlers.PlaceOrder)             // Place a new order - completed
		orderRoutes.GET("/:id", handlers.GetOrder)             // Fetch order details - completed
		orderRoutes.GET("/", handlers.GetOrdersWithPagination) // Fetch order details - completed, user preload not coming
		orderRoutes.PATCH("/:id/cancel", handlers.CancelOrder) // Cancel an order - completed
		orderRoutes.GET("/:id/history", handlers.HistoryOrder) // Cancel an order
	}
}
