package routes

import (
	"ordermanagement/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes initializes user-related routes
func SetupUserRoutes(api *gin.RouterGroup) {
	userRoutes := api.Group("/users") // Base path for user routes
	{
		userRoutes.POST("/", handlers.CreateUser) // Create a user
		// userRoutes.GET("/:id", handlers.GetUser)      // Get user by ID
		// userRoutes.GET("/", handlers.GetUsers)        // Get all users
		// userRoutes.PUT("/:id", handlers.UpdateUser)   // Update user
		// userRoutes.DELETE("/:id", handlers.DeleteUser) // Delete user
	}
}
