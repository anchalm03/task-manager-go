package services

import (
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func TokenHandler(c *gin.Context) {
	tokenString, err := middlewares.GenerateToken("user123")
	if err != nil {
		c.JSON(500, gin.H{"error": "could not sign token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString})
}
