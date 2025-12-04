package middlewares

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT token for a given userID
// can be extended to phone, other user params
func GenerateToken(userID string, role string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtSecret)
}
