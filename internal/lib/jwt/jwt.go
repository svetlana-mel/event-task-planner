package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

var (
	ErrTokenExpired = errors.New("token expired")
)

func NewToken(userID uint64, email string, jwtSecret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userID
	claims["email"] = email
	claims["expiration"] = time.Now().Add(duration).Unix()

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(token *jwt.Token) (*models.User, error) {
	const op = "lib.jwt.ValidateToken"

	if !token.Valid {
		return nil, fmt.Errorf("%s: token invalid", op)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: token error get claims", op)
	}

	uid := claims["uid"].(string)
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	email := claims["email"].(string)

	exp := claims["expiration"].(string)
	expirationTime, err := time.Parse(time.UnixDate, exp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if expirationTime.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return &models.User{
		UserID: userID,
		Email:  email,
	}, nil
}
