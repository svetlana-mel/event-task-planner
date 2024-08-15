package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
	base "github.com/svetlana-mel/event-task-planner/internal/repository"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres/converter"
)

func (r *repository) GetTask(ctx context.Context, taskID uint64) (*models.Task, error) {
	op := "repository.postgres.GetTask"
	row := r.pool.QueryRow(ctx, `select * from task where task_id=$1`, taskID)

	task, err := converter.TaskRowToModel(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, base.ErrTaskNotExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (r *repository) CreateTask(ctx context.Context, task models.Task) error {
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

func (r *repository) GetAllTasks(ctx context.Context, status string) ([]models.Task, error) {
	op := "repository.postgres.GetAllTask"

	sql := `select * from task`
	switch status {
	case "active":
		sql = `select * from task where completed_at is null`
	case "completed":
		sql = `select * from task where completed_at is not null`
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
