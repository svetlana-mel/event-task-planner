package models

import (
	"time"
)

type Event struct {
	EventID  uint64 `db:"event_id" json:"event_id"`
	FkUserID uint64 `db:"fk_user_id" json:"fk_user_id"`

	// editible by user fields
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	DateTime    time.Time `db:"date_time" json:"date_time"`
	Tasks       []Task    `db:"tasks" json:"tasks"`

	// pointer types uses for fields that can be NULL
	CanceledAt *time.Time `db:"canceled_at" json:"canceled_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}
