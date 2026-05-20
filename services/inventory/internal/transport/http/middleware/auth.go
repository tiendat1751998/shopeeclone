package middleware

import (
	"github.com/gin-gonic/gin"
	sharedMiddleware "github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/services/inventory/internal/config"
)

// JWTAuth validates JWT tokens and extracts user context.
// Delegates to the shared middleware implementation.
func JWTAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return sharedMiddleware.JWTAuth(cfg.AccessSecret)
}

// RequireRole creates middleware that checks if the authenticated user has the required role.
// Delegates to the shared middleware implementation.
func RequireRole(roles ...string) gin.HandlerFunc {
	return sharedMiddleware.RequireRole(roles...)
}
