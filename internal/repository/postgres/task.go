package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
	base "github.com/svetlana-mel/event-task-planner/internal/repository"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres/converter"
)

func (r *repository) GetTask(ctx context.Context, taskID uint64) (*models.Task, error) {
	op := "repository.postgres.GetTask"
	rows, err := r.pool.Query(ctx, `select * from task where task_id=$1`, taskID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, base.ErrTaskNotExists)
	}
	defer rows.Close()

	tasks, err := converter.TaskRowsToModel(rows)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("%s: %w", op, base.ErrTaskNotExists)
	}

	return &tasks[0], nil
}

func (r *repository) CreateTask(ctx context.Context, task *models.Task) error {
	op := "repository.postgres.CreateTask"

	_, err := r.pool.Exec(ctx,
		`insert into "task" 
		(name, description, list, start_date_time, 
		end_date_time, fk_event_id, fk_user_id)
		values 
		(	$1, $2, $3, $4,
			nullif($5, to_timestamp(0)), 
			nullif($6, 0), $7
		)`,
		task.Name,
		task.Description,
		task.List,
		task.StartDateTime,
		task.EndDateTime,
		task.FkEventID,
		task.FkUserID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) UpdateTask(ctx context.Context, task *models.Task) error {
	// function to update content fields:
	// 		name, description, list
	// and update start-end timestamps

	op := "repository.postgres.UpdateTask"

	sql := `update "task" set
	name = $1,
	description = $2,
	list = $3,
	start_date_time = $4,
	end_date_time = $5
	where task_id = $6
	`
	_, err := r.pool.Exec(ctx, sql,
		task.Name,
		task.Description,
		task.List,
		task.StartDateTime,
		task.EndDateTime,
		task.TaskID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) SetTaskCompletionStatus(ctx context.Context, taskID uint64, completed bool) error {
	// function to update task completed_at field
	// if completed_at == NULL (not completed), sets completed_at as timestamp (completed)
	// if completed_at NOT NULL (completed), sets completed_at to NULL (not completed)

	op := "repository.postgres.SetTaskCompletionStatus"

	sql := `update "task" set completed_at = $1 where task_id = $2`

	var completedAt *time.Time
	if completed {
		t := time.Now()
		completedAt = &t
	}

	_, err := r.pool.Exec(ctx, sql, completedAt, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) DeleteTask(ctx context.Context, taskID uint64) error {
	op := "repository.postgres.DeleteTask"

	sql := `update "task" set deleted_at = $1 where task_id = $2`

	timestamp := time.Now()

	_, err := r.pool.Exec(ctx, sql, timestamp, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) GetAllTasks(ctx context.Context, status string) ([]models.Task, error) {
	op := "repository.postgres.GetAllTasks"

	sql := `select * from task`
	switch status {
	case "active":
		sql = `select * from task where completed_at is null`
	case "completed":
		sql = `select * from task where completed_at is not null`
	case "unassigned":
		sql = `select * from task where event_id is null`
	}

	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, base.ErrTaskNotExists)
	}
	defer rows.Close()

	tasks, err := converter.TaskRowsToModel(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, base.ErrTaskNotExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
