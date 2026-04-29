package services

import (
	"task_manager/middlewares"
	"task_manager/services.go/auth"
	"task_manager/services.go/comments"
	"task_manager/services.go/projects"
	"task_manager/services.go/tasks"
	"task_manager/services.go/users"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/ping", HealthCheck)

	r.POST("/register", auth.RegisterUser)
	r.POST("/login", auth.LoginUser)

	authentication := r.Group("/api")
	authentication.Use(middlewares.JWTAuthMiddleware())
	{
		authentication.GET("/me", users.GetMe)

		authentication.POST("/projects", projects.CreateProject)
		authentication.GET("/projects", projects.ListProjects)
		authentication.GET("/projects/:id", projects.GetProject)
		authentication.PATCH("/projects/:id", projects.UpdateProject)
		authentication.DELETE("/projects/:id", projects.DeleteProject)

		authentication.POST("/projects/:id/tasks", tasks.CreateTask)
		authentication.GET("/projects/:id/tasks", tasks.ListTasksByProject)

		authentication.GET("/me/tasks", tasks.ListMyTasks)
		authentication.PATCH("/tasks/:id", tasks.UpdateTask)
		authentication.DELETE("/tasks/:id", tasks.DeleteTask)

		authentication.POST("/tasks/:id/comments", comments.CreateComment)
		authentication.GET("/tasks/:id/comments", comments.ListCommentsByTask)
		authentication.DELETE("/comments/:id", comments.DeleteComment)
	}
}
