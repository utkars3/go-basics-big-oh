package middleware

import (
	"log"
	"net/http"
	"ordermanagement/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve userID from context (set by JWT middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No User ID found"})
			c.Abort()
			return
		}

		// Query PostgreSQL to check user role using GORM's Raw SQL query
		var role string
		err := config.DBPostgreSQL.Raw("SELECT role FROM users WHERE id = ?", userID).Scan(&role).Error
		if err != nil {
			// Handle the case where no rows are returned
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User not found"})
			} else {
				// Log and handle other errors
				log.Println("Error fetching user role:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
			c.Abort()
			return
		}

		// Allow only if the role is "admin"
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Only admin can access this"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
