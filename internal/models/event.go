package models

import "time"

type Event struct {
	EventID  uint64
	Name     string
	DateTime time.Time
	FkUserID uint64

	// pointer types uses for fields that can be NULL
	CanceledAt *time.Time
	DeletedAt  *time.Time
}
