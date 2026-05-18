package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey string

const (
	CorrelationIDKey ContextKey = "correlation_id"
	RequestIDKey     ContextKey = "request_id"
	UserIDKey        ContextKey = "user_id"
	UserRolesKey     ContextKey = "user_roles"
	DeviceInfoKey    ContextKey = "device_info"
)

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		c.Set(string(CorrelationIDKey), correlationID)
		c.Header("X-Correlation-ID", correlationID)

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(string(RequestIDKey), requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
