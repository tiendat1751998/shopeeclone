package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RedisSlidingWindowLimiter implements a sliding-window rate limiter using Redis sorted sets.
func RedisSlidingWindowLimiter(redisClient *redis.Client, maxRequests int, window time.Duration, keyPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:%s:%s", keyPrefix, c.ClientIP())
		now := time.Now().UnixNano()
		windowStart := now - window.Nanoseconds()

		pipe := redisClient.Pipeline()
		pipe.ZRemRangeByScore(c.Request.Context(), key, "0", fmt.Sprintf("%d", windowStart))
		pipe.ZCard(c.Request.Context(), key)
		pipe.ZAdd(c.Request.Context(), key, redis.Z{Score: float64(now), Member: now})
		pipe.Expire(c.Request.Context(), key, window)

		cmds, err := pipe.Exec(c.Request.Context())
		if err != nil {
			c.Next()
			return
		}

		count := cmds[1].(*redis.IntCmd).Val()
		if int(count) >= maxRequests {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error_code": "RATE_LIMIT_EXCEEDED",
				"message":    fmt.Sprintf("too many requests, max %d per %s", maxRequests, window),
			})
			return
		}

		c.Next()
	}
}
