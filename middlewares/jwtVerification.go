package middlewares

import (
	"os"
	"strings"
	"task_manager/errorcodes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware validates JWT tokens in Authorization header.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorcode := errorcodes.NoError

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorcode = errorcodes.BadRequest
			c.JSON(errorcode.HttpStatusCode(), gin.H{
				"success":       false,
				"error_message": "missing Authorization header",
				"error_code":    errorcode,
			})
			c.Abort()
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			errorcode = errorcodes.BadRequest
			c.JSON(errorcode.HttpStatusCode(), gin.H{
				"success":       false,
				"error_message": "invalid Authorization header format",
				"error_code":    errorcode,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			errorcode = errorcodes.BadRequest
			c.JSON(errorcode.HttpStatusCode(), gin.H{
				"success":       false,
				"error_message": "invalid or expired token",
				"error_code":    errorcode,
			})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["sub"])
			c.Set("role", claims["role"])
		} else {
			errorcode = errorcodes.InternalServerError
			c.JSON(errorcode.HttpStatusCode(), gin.H{
				"success":       false,
				"error_message": errorcode.Message(),
				"error_code":    errorcode,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
