package models

import "time"

type User struct {
	userID             uint64
	name               string
	email              string
	password           string
	createdDateTime    time.Time
	updatedDateTime    time.Time
	lastLogin          time.Time
	refreshToken       string
	refreshTokenExpiry time.Time
}
