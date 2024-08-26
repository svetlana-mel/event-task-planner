package repository

import (
	"context"
	"errors"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

var (
	ErrTaskNotExists     = errors.New("task not exists")
	ErrEventNotExists    = errors.New("event not exists")
	ErrUserNotExists     = errors.New("user not exists")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type PlannerRepository interface {
	// tasks crud
	GetTask(ctx context.Context, taskID uint64) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	SetTaskCompletionStatus(ctx context.Context, taskID uint64, completed bool) error
	DeleteTask(ctx context.Context, taskID uint64) error
	GetAllTasks(ctx context.Context, status string) ([]models.Task, error)

	// events crud
	GetEvent(ctx context.Context, eventID uint64) (*models.Event, error)
	CreateEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, event *models.Event) error
	SetEventCanceledStatus(ctx context.Context, eventID uint64, canceled bool) error
	DeleteEvent(ctx context.Context, eventID uint64) error
	GetAllEvents(ctx context.Context, status string) ([]models.Event, error)

	// users crud
	CreateUser(ctx context.Context, username string, email string, passHash []byte) (uint64, error)
	GetUser(ctx context.Context, email string) (*models.User, error)

	// methods to manage event placeholder for tasks without a referenced event
	GetDefaultEventID(ctx context.Context) uint64
	GetDefaultUserID(ctx context.Context) uint64

	// dev tools
	CreateTmpUser(ctx context.Context) (uint64, error)

	// free resources if needed
	Close()
}
