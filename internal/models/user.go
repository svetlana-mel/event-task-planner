package models

import "time"

type User struct {
	UserID             uint64    `db:"user_id"`
	Name               string    `db:"name"`
	Email              string    `db:"email"`
	Password           string    `db:"password"`
	CreatedDateTime    time.Time `db:"created_date_time"`
	UpdatedDateTime    time.Time `db:"updated_date_time"`
	LastLogin          time.Time `db:"last_login"`
	RefreshToken       string    `db:"refresh_token"`
	RefreshTokenExpiry time.Time `db:"refresh_token_expiry"`
}
