package models

import "time"

type Event struct {
	eventID  uint64
	name     string
	dateTime time.Time
	canceled bool
	deleted  bool
	fkUserID uint64
}
