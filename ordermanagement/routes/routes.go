package routes

import "github.com/gin-gonic/gin"

// SetupRoutes initializes all routes
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api") // Base API path

	SetupUserRoutes(api) // Mount user routes at /api/users
	SetupOrderRoutes(api)
	SetupInventoryRoutes(api)
}
