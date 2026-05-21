package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GinJWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(bearer, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		claims, err := ParseToken(parts[1], secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		userID, _ := claims["sub"].(string)
		if userID == "" {
			userID, _ = claims["user_id"].(string)
		}
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing user identifier"})
			return
		}
		role := "buyer"
		if r, ok := claims["role"].(string); ok && r != "" {
			role = r
		} else if roles, ok := claims["roles"].([]interface{}); ok && len(roles) > 0 {
			if r, ok := roles[0].(string); ok {
				role = r
			}
		}
		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}
