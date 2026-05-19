package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/gateway/internal/config"
)

func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	allowMethods := strings.Join(cfg.AllowedMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowedHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposedHeaders, ", ")
	maxAge := strconv.Itoa(int(cfg.MaxAge.Seconds()))

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" {
			allowedOrigin := ""

			// [SECURITY] Check len before accessing index 0
			if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] == "*" {
				if !cfg.AllowCredentials {
					allowedOrigin = "*"
				} else {
					// [SECURITY] When credentials are used, echo the specific origin
					allowedOrigin = origin
				}
			} else {
				for _, allowed := range cfg.AllowedOrigins {
					if allowed == origin {
						allowedOrigin = origin
						break
					}
				}
			}

			if allowedOrigin != "" {
				c.Header("Access-Control-Allow-Origin", allowedOrigin)
			}

			c.Header("Access-Control-Allow-Methods", allowMethods)
			c.Header("Access-Control-Allow-Headers", allowHeaders)
			c.Header("Access-Control-Expose-Headers", exposeHeaders)
			c.Header("Access-Control-Max-Age", maxAge)

			if cfg.AllowCredentials && allowedOrigin != "" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
