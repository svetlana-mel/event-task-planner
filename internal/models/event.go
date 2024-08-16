package models

import "time"

type Event struct {
	EventID  uint64 `db:"event_id"`
	FkUserID uint64 `db:"fk_user_id"`

	// editible by user fields
	Name        string    `db:"name"`
	Description string    `db:"description"`
	DateTime    time.Time `db:"date_time"`

	// pointer types uses for fields that can be NULL
	CanceledAt *time.Time `db:"canceled_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}
