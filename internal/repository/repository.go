package repository

import (
	"context"
	"errors"
	"time"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

var (
	ErrTaskNotExists = errors.New("task not exists")
)

type PlannerRepository interface {
	GetTask(ctx context.Context, taskID uint64) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	SetTaskCompletionStatus(ctx context.Context, taskID uint64, completed bool) error
	DeleteTask(ctx context.Context, taskID uint64) error
	GetAllTasks(ctx context.Context, completed bool) ([]models.Task, error)

	GetEvent(ctx context.Context, eventID uint64) (models.Event, error)
	CreateEvent(ctx context.Context, event models.Event) error
	UpdateEvent(ctx context.Context, event models.Event) error
	GetAllEvents(ctx context.Context, canceled bool) ([]models.Event, error)

	GetEventTasks(ctx context.Context, eventID uint64) ([]models.Task, error)
	GetTodayTasks(ctx context.Context, date time.Time) ([]models.Task, error)
	GetTodayEvents(ctx context.Context, date time.Time) ([]models.Event, error)

	// dev tools
	CreateTmpUser(ctx context.Context) (uint64, error)

	// free resources if needed
	Close()
}
