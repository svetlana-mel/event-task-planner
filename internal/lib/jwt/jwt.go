package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

func NewToken(user *models.User, jwtSecret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.UserID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
