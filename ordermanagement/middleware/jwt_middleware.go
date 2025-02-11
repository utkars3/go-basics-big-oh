package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Secret key used for signing JWT tokens
// var jwtSecret = []byte("your_secret_key")
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims structure for JWT payload
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWTAuthMiddleware verifies JWT in Authorization header
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Check if token is provided and has "Bearer " prefix
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		// Extract actual token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		parsedUserID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		userIDFromHeader := c.GetHeader("user_id")

		// fmt.Println("------------------------", userIDFromHeader)
		// fmt.Println("------------------------", claims.UserID)
		if userIDFromHeader == "" || userIDFromHeader != claims.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID mismatch"})
			c.Abort()
			return
		}

		// Store user ID in context for later use
		c.Set("userID", parsedUserID)
		c.Next()
	}
}
