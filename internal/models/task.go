package models

import "time"

type Task struct {
	TaskID   uint64 `db:"task_id"`
	FkUserID uint64 `db:"fk_user_id"`

	// editible by user fields
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	List          []string  `db:"list"`
	StartDateTime time.Time `db:"start_date_time"`

	// pointer data types for fields that can be NULL
	// fo example in CompletedAt case
	// the NULL type interpreted as not completed
	// and not NULL(with actual date) as completed
	EndDateTime *time.Time `db:"end_date_time"`
	FkEventID   *uint64    `db:"fk_event_id"`
	CompletedAt *time.Time `db:"completed_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
