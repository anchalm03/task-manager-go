package services

import (
	"log"
	"task_manager/errorcodes"
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

type TokenHandlerResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         string               `json:"data,omitempty"`
}

func TokenHandler(c *gin.Context) {
	error := errorcodes.NoError
	tokenString, err := middlewares.GenerateToken("user123")
	if err != nil {
		log.Println(`error: could not sign token`)
		error = errorcodes.InternalServerError
		c.JSON(error.HttpStatusCode(), TokenHandlerResponse{Success: false, ErrorCode: error, ErrorMessage: error.Message()})
		return
	}
	c.JSON(error.HttpStatusCode(), TokenHandlerResponse{Success: true, ErrorCode: error, ErrorMessage: error.Message(), Data: tokenString})
}
