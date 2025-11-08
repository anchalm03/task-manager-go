package services

import (
	"task_manager/errorcodes"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         string               `json:"data"`
}

// HealthCheck returns a simple health status
func HealthCheck(c *gin.Context) {
	errorcode := errorcodes.NoError
	c.JSON(errorcode.HttpStatusCode(), HealthResponse{Success: true, ErrorCode: errorcode, ErrorMessage: errorcode.Message(), Data: "pong"})
}
