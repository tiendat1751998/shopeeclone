package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutHandler(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{}, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					c.Error(fmt.Errorf("panic in request handler: %v", r))
				}
				select {
				case done <- struct{}{}:
				default:
				}
			}()
			c.Next()
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error_code": "GATEWAY_TIMEOUT",
				"message":    "request timed out",
			})
		}
	}
}
