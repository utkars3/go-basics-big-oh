package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ip := c.Request.Header.Get("X-Forwarded-For")
		// if ip == "" {
		// 	ip = c.ClientIP() // If not available, fall back to ClientIP
		// } else {
		// 	ip = strings.Split(ip, ",")[0] // Get the first IP in the list
		// }

		// fmt.Println("Client IP:", ip) // Log the client IP

		c.Header("Access-Control-Allow-Origin", "https://localhost:3005") // Allow only specific domains
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		// fmt.Println("logging -----------------------", c.ClientIP())

		// Handle preflight OPTIONS request
		// A preflight request is a request automatically sent by the browser before making an actual cross-origin request.
		// It is sent to check whether the actual request is allowed by the server.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			// fmt.Println("logging -----------------------", c.ClientIP())

			return
		}

		c.Next()
	}
}
