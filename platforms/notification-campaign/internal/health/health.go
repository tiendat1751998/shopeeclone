package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Checker struct {
	appName string
	version string
}

func NewChecker(appName, version string) *Checker {
	return &Checker{appName: appName, version: version}
}

func (h *Checker) LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive", "app": h.appName})
	}
}

func (h *Checker) ReadinessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready", "app": h.appName})
	}
}
