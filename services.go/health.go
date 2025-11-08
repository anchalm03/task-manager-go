package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns a simple health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "pong",
	})
}
