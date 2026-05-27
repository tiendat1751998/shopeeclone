package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth returns a gin middleware that validates the X-API-Key header
// against the GLOBAL_INFRA_API_KEY environment variable.
func APIKeyAuth() gin.HandlerFunc {
	expectedKey := os.Getenv("GLOBAL_INFRA_API_KEY")
	return func(c *gin.Context) {
		if expectedKey == "" {
			c.Next()
			return
		}
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code": "MISSING_API_KEY",
				"message":    "X-API-Key header is required",
			})
			return
		}
		if key != expectedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code": "INVALID_API_KEY",
				"message":    "invalid API key",
			})
			return
		}
		c.Next()
	}
}
