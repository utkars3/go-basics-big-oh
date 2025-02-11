package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Log request details
		duration := time.Since(startTime)
		log.Printf("[%s] %s %d %s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
