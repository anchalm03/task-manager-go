package users

import (
	"errors"
	"task_manager/errorcodes"
	"task_manager/middlewares"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// ---- view & response shapes ----

type UserView struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func toView(u *models.User) UserView {
	return UserView{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}
}

type UserResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    errorcodes.ErrorCode `json:"errorcode"`
	ErrorMessage string               `json:"error"`
	Data         *UserView            `json:"data,omitempty"`
}

// ---- helpers ----

func msgFor(code errorcodes.ErrorCode, override string) string {
	if override != "" {
		return override
	}
	return code.Message()
}

func writeUser(c *gin.Context, code errorcodes.ErrorCode, override string, data *UserView) {
	c.JSON(code.HttpStatusCode(), UserResponse{
		Success:      code == errorcodes.NoError,
		ErrorCode:    code,
		ErrorMessage: msgFor(code, override),
		Data:         data,
	})
}

// ---- handlers ----

// GET /api/me — returns the caller's profile (no password hash).
func GetMe(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		writeUser(c, errorcodes.Unauthorized, "", nil)
		return
	}

	user, err := models.GetUserById(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeUser(c, errorcodes.NotFound, "", nil)
			return
		}
		writeUser(c, errorcodes.InternalServerError, "", nil)
		return
	}

	view := toView(user)
	writeUser(c, errorcodes.NoError, "", &view)
}
