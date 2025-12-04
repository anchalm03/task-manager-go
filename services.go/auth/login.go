package auth

import (
	"task_manager/db"
	"task_manager/errorcodes"
	"task_manager/middlewares"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Token        string               `json:"token,omitempty"`
}

func LoginUser(c *gin.Context) {
	var req LoginRequest
	errorcode := errorcodes.NoError

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		errorcode = errorcodes.BadRequest
		c.JSON(errorcode.HttpStatusCode(), LoginResponse{Success: false, ErrorCode: errorcode, ErrorMessage: errorcode.Message()})
		return
	}

	// Find user by email
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		errorcode = errorcodes.NotFound
		c.JSON(errorcode.HttpStatusCode(), LoginResponse{Success: false, ErrorCode: errorcode, ErrorMessage: "user not found"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		errorcode = errorcodes.Unauthorized
		c.JSON(errorcode.HttpStatusCode(), LoginResponse{Success: false, ErrorCode: errorcode, ErrorMessage: "invalid credentials"})
		return
	}

	// Generate JWT
	token, err := middlewares.GenerateToken(user.Id.String(), user.Role)
	if err != nil {
		errorcode = errorcodes.InternalServerError
		c.JSON(errorcode.HttpStatusCode(), LoginResponse{Success: false, ErrorCode: errorcode, ErrorMessage: errorcode.Message()})
		return
	}

	c.JSON(errorcode.HttpStatusCode(), LoginResponse{Success: true, ErrorCode: errorcode, ErrorMessage: errorcode.Message(), Token: token})
}
