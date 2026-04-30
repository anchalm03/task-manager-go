package auth

import (
	"task_manager/errorcodes"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// you can integrate googleAPI also for login

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	errorcode := errorcodes.NoError
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode = errorcodes.BadRequest
		c.JSON(errorcode.HttpStatusCode(), RegisterResponse{Success: false, ErrorMessage: errorcode.Message(), ErrorCode: errorcode})
		return
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "member",
	}

	if err := models.CreateUser(&user); err != nil {
		errorcode = errorcodes.InternalServerError
		c.JSON(errorcode.HttpStatusCode(), RegisterResponse{Success: false, ErrorMessage: errorcode.Message(), ErrorCode: errorcode})
		return
	}

	c.JSON(errorcode.HttpStatusCode(), RegisterResponse{Success: true, ErrorMessage: errorcode.Message(), ErrorCode: errorcode})
}
