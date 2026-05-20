package middleware

import (
	"github.com/gin-gonic/gin"
	sharedMiddleware "github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/services/shipment/internal/config"
)

// JWTAuth validates JWT tokens and extracts user context.
// Delegates to the standardized shared middleware implementation.
func JWTAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return sharedMiddleware.JWTAuth(cfg.AccessSecret)
}
