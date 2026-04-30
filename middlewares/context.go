package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// GetUserID pulls the authenticated user's ID out of the Gin context.
// Returns an error if the middleware hasn't run or the value is malformed.
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	raw, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, errors.New("user_id not set in context")
	}
	idStr, ok := raw.(string)
	if !ok {
		return uuid.Nil, errors.New("user_id in context is not a string")
	}
	return uuid.FromString(idStr)
}

// GetRole returns the caller's role claim ("member", "admin", etc).
func GetRole(c *gin.Context) string {
	raw, ok := c.Get("role")
	if !ok {
		return ""
	}
	role, _ := raw.(string)
	return role
}
