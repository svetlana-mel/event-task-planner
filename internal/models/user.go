package models

import "time"

type User struct {
	UserID          uint64    `db:"user_id"`
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	PassHash        []byte    `db:"pass_hash"`
	CreatedDateTime time.Time `db:"created_date_time"`
}
