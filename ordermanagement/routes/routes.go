package routes

import (
	"ordermanagement/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all routes
func SetupRoutes(r *gin.Engine) {
	rl := middleware.NewRateLimiter(1, 2)
	api := r.Group("/api") // Base API path
	api.Use(middleware.LoggerMiddleware())
	api.Use(rl.RateLimitMiddleware())

	
	SetupUserRoutes(api) // Mount user routes at /api/users
	SetupOrderRoutes(api)
	SetupInventoryRoutes(api)
}
