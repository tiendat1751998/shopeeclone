package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"go.uber.org/zap"
)

type SecurityConfig struct {
	MaxBodySize    int64
	AllowedHosts   []string
	TrustedProxies []string
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")

		c.Next()
	}
}

func BodySizeLimiter(maxBodySize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		}
		c.Next()
	}
}

func RequestSanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		for key, values := range query {
			for i, v := range values {
				values[i] = sanitizeInput(v)
			}
			query[key] = values
		}
		c.Request.URL.RawQuery = query.Encode()

		c.Next()
	}
}

func sanitizeInput(input string) string {
	input = strings.ReplaceAll(input, "<script", "&lt;script")
	input = strings.ReplaceAll(input, "</script>", "&lt;/script&gt;")
	input = strings.ReplaceAll(input, "javascript:", "")
	input = strings.ReplaceAll(input, "onerror=", "onerror-disabled=")
	input = strings.ReplaceAll(input, "onload=", "onload-disabled=")
	return input
}

func IPThrottler(trustedProxies []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, proxy := range trustedProxies {
			if strings.HasPrefix(clientIP, proxy) {
				c.Next()
				return
			}
		}
		c.Next()
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		bodySize := c.Writer.Size()
		userID, _ := c.Get(string(UserIDKey))

		log := observability.LogWithTrace(c.Request.Context())
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
			zap.Int("body_size", bodySize),
		}
		if userID != nil {
			fields = append(fields, zap.String("user_id", fmt.Sprintf("%v", userID)))
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				fields = append(fields, zap.String("error", e.Err.Error()))
			}
			log.Error("request completed with errors", fields...)
		} else if status >= 500 {
			log.Error("server error", fields...)
		} else if status >= 400 {
			log.Warn("client error", fields...)
		} else {
			log.Info("request completed", fields...)
		}
	}
}

func AntiAbuse() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_code": "MISSING_USER_AGENT",
				"message":    "User-Agent header is required",
			})
			return
		}
		if len(userAgent) > 500 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_code": "INVALID_USER_AGENT",
				"message":    "User-Agent header too long",
			})
			return
		}

		contentType := c.GetHeader("Content-Type")
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			if contentType == "" {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error_code": "MISSING_CONTENT_TYPE",
					"message":    "Content-Type header is required for write operations",
				})
				return
			}
		}

		query := c.Request.URL.RawQuery
		if len(query) > 2000 {
			c.AbortWithStatusJSON(http.StatusRequestURITooLong, gin.H{
				"error_code": "QUERY_TOO_LONG",
				"message":    "Query string exceeds maximum length",
			})
			return
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			c.Header("Content-Encoding", "gzip")
		}

		c.Next()
	}
}

func DeviceMetadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceInfo := map[string]string{
			"user_agent":      c.GetHeader("User-Agent"),
			"device_id":       c.GetHeader("X-Device-ID"),
			"device_type":     c.GetHeader("X-Device-Type"),
			"platform":        c.GetHeader("X-Platform"),
			"app_version":     c.GetHeader("X-App-Version"),
			"accept_language": c.GetHeader("Accept-Language"),
		}
		c.Set(string(DeviceInfoKey), deviceInfo)
		c.Next()
	}
}

func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateRequest(c); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_code": "VALIDATION_ERROR",
				"message":    err.Error(),
			})
			return
		}
		c.Next()
	}
}

func validateRequest(c *gin.Context) error {
	method := c.Request.Method
	path := c.Request.URL.Path

	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal detected")
	}

	if method == http.MethodTrace {
		return fmt.Errorf("TRACE method not allowed")
	}

	if method == http.MethodConnect {
		return fmt.Errorf("CONNECT method not allowed")
	}

	host := c.GetHeader("Host")
	if host == "" {
		return fmt.Errorf("Host header is required")
	}

	contentLength := c.GetHeader("Content-Length")
	if contentLength != "" {
		length, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid Content-Length header")
		}
		if length < 0 {
			return fmt.Errorf("negative Content-Length")
		}
	}

	return nil
}
