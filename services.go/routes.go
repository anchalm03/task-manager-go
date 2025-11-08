package services

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/ping", HealthCheck)
	r.GET("/token", TokenHandler)
}
