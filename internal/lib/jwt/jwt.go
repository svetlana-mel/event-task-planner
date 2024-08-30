package jwt

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("token invalid")
)

func NewToken(userID uint64, email string, jwtPrivateKey *ecdsa.PrivateKey, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userID
	claims["email"] = email
	claims["expiration"] = time.Now().Add(duration).Unix()

	tokenStr, err := token.SignedString(jwtPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(tokenString string, jwtPublicKey *ecdsa.PublicKey) (*models.User, error) {
	const op = "lib.jwt.ValidateToken"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		signinMethod := token.Method
		if _, ok := signinMethod.(*jwt.SigningMethodECDSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return jwtPublicKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
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
