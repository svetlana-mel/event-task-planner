package repository

import (
	"time"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

type PlannerRepository interface {
	GetTask(taskID uint64) (models.Task, error)
	CreateTask(task models.Task) (uint64, error)
	UpdateTask(task models.Task) error
	GetAllTasks(completed bool) ([]models.Task, error)

	GetEvent(eventID uint64) (models.Event, error)
	CreateEvent(event models.Event) (uint64, error)
	UpdateEvent(event models.Event) error
	GetAllEvents(canceled bool) ([]models.Event, error)

	GetEventTasks(eventID uint64) ([]models.Task, error)
	GetTodayTasks(date time.Time) ([]models.Task, error)
	GetTodayEvents(date time.Time) ([]models.Event, error)
}
