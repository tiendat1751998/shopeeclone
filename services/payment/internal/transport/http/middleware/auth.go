package middleware

import (
	"github.com/gin-gonic/gin"
	authpkg "github.com/shopee-clone/shopee/packages/go-shared/pkg/auth"
	"github.com/shopee-clone/shopee/services/payment/internal/config"
)

func JWTAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return authpkg.GinJWTAuth(cfg.AccessSecret)
}
