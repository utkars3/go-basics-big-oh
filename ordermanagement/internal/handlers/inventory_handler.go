package handlers

import (
	"net/http"
	"ordermanagement/internal/models"
	"ordermanagement/internal/services"

	"github.com/gin-gonic/gin"
)

func CreateInventory(c *gin.Context) {
	var item models.Inventory
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to create or update inventory
	createdItem, err := services.CreateInventoryService(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory processed successfully", "data": createdItem})
}

func GetInventoryHandler(c *gin.Context) {
	inventory, err := services.GetAllInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inventory": inventory})
}
