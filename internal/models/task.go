package models

import "time"

type Task struct {
	taskID   uint64
	fkUserID uint64

	// editible fields
	name          string
	description   string
	list          []string
	startDateTime time.Time
	endDateTime   time.Time
	fkEventID     uint64
	completed     bool
	deleted       bool
}
