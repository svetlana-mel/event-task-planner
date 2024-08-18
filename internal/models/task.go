package models

import (
	"time"
)

type Task struct {
	TaskID   uint64 `db:"task_id" json:"task_id"`
	FkUserID uint64 `db:"fk_user_id" json:"fk_user_id"`

	// editible by user fields
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	List          []string  `db:"list" json:"list"`
	StartDateTime time.Time `db:"start_date_time" json:"start_date_time"`

	// pointer data types for fields that can be NULL
	// fo example in CompletedAt case
	// the NULL type interpreted as not completed
	// and not NULL(with actual date) as completed
	EndDateTime *time.Time `db:"end_date_time" json:"end_date_time"`
	FkEventID   *uint64    `db:"fk_event_id" json:"fk_event_id"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}
