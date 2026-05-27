package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		latency := time.Since(start)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
		}
		if requestID, exists := c.Get("request_id"); exists {
			if rid, ok := requestID.(string); ok {
				fields = append(fields, zap.String("request_id", rid))
			}
		}
		if c.Writer.Status() >= 500 {
			zap.L().Error("request failed", fields...)
		} else {
			zap.L().Info("request completed", fields...)
		}
	}
}
