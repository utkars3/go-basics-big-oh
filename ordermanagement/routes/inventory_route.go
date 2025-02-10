package routes

import (
	"ordermanagement/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupInventoryRoutes initializes inventory-related routes
func SetupInventoryRoutes(api *gin.RouterGroup) {
	inventoryRoutes := api.Group("/inventory") // Base path for inventory routes
	{
		inventoryRoutes.POST("/", handlers.CreateInventory)    // Add new inventory item
		inventoryRoutes.GET("/", handlers.GetInventoryHandler) 
	}
}
