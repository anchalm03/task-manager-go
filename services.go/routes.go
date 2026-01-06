package services

import (
	"task_manager/middlewares"
	"task_manager/services.go/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/ping", HealthCheck)
	r.GET("/token", TokenHandler)

	r.POST("/register", auth.RegisterUser)
	r.POST("/login", auth.LoginUser)

	authentication := r.Group("/api")
	authentication.Use(middlewares.JWTAuthMiddleware())
	{
		// all the APIs that need auth
	}
}
