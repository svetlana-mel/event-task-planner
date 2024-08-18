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

func (r *repository) GetEvent(ctx context.Context, eventID uint64) (*models.Event, error) {
	op := "repository.postgres.GetEvent"

	sql := `
	SELECT e.*, COALESCE(json_agg(
               json_build_object(
                   'task_id', t.task_id, 
                   'name', t.name,
                   'description', t.description,
                   'list', t.list,
                   'start_date_time', t.start_date_time,
                   'end_date_time', t.end_date_time,
                   'fk_event_id', t.fk_event_id,
                   'completed_at', t.completed_at,
                   'deleted_at', t.deleted_at,
                   'fk_user_id', t.fk_user_id
               )
           ), '[]') as tasks
    FROM "event" as e
    LEFT JOIN "task" as t ON e.event_id = t.fk_event_id
	WHERE e.event_id=$1
    GROUP BY e.event_id

	`
	rows, _ := r.pool.Query(ctx, sql, eventID)
	event, err := pgx.CollectExactlyOneRow(rows, converter.EventRowToModel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &event, nil
}

func (r *repository) CreateEvent(ctx context.Context, event *models.Event) error {
	op := "repository.postgres.CreateEvent"

	_, err := r.pool.Exec(ctx,
		`insert into "event" 
		(name, description, date_time, fk_user_id)
		values 
		($1, $2, $3, $4)`,
		event.Name,
		event.Description,
		event.DateTime,
		event.FkUserID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) UpdateEvent(ctx context.Context, event *models.Event) error {
	// function to update content fields:
	// 		name, description, date_time

	op := "repository.postgres.UpdateEvent"

	sql := `update "event" set
	name = $1,
	description = $2,
	date_time = $3
	where event_id = $4
	`
	_, err := r.pool.Exec(ctx, sql,
		event.Name,
		event.Description,
		event.DateTime,
		event.EventID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) SetEventCanceledStatus(ctx context.Context, eventID uint64, canceled bool) error {
	op := "repository.postgres.SetEventCanceledStatus"

	sql := `update "event" set
	canceled_at = $1
	where event_id = $2
	`

	var canceledAt *time.Time
	if canceled {
		t := time.Now()
		canceledAt = &t
	}

	_, err := r.pool.Exec(ctx, sql, canceledAt, eventID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) DeleteEvent(ctx context.Context, eventID uint64) error {
	op := "repository.postgres.DeleteEvent"

	sql := `update "event" set deleted_at = $1 where event_id = $2`

	timestamp := time.Now()

	_, err := r.pool.Exec(ctx, sql, timestamp, eventID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) GetAllEvents(ctx context.Context, status string) ([]models.Event, error) {
	op := "repository.postgres.GetAllEvents"

	completedFilter := ""

	switch status {
	case "active":
		completedFilter = "AND canceled_at IS NULL"
	case "completed":
		completedFilter = "AND canceled_at IS NOT NULL"
	}

	sql := fmt.Sprintf(`
	SELECT e.*, COALESCE(json_agg(
			json_build_object(
				'task_id', t.task_id, 
				'name', t.name,
				'description', t.description,
				'list', t.list,
				'start_date_time', t.start_date_time,
				'end_date_time', t.end_date_time,
				'fk_event_id', t.fk_event_id,
				'completed_at', t.completed_at,
				'deleted_at', t.deleted_at,
				'fk_user_id', t.fk_user_id
			)
		), '[]') as tasks
    FROM "event" as e
    LEFT JOIN "task" as t ON e.event_id = t.fk_event_id
	WHERE e.deleted_at IS NULL %s
    GROUP BY e.event_id`, completedFilter)

	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, base.ErrEventNotExists)
	}
	defer rows.Close()

	events, err := converter.EventsRowsToModel(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, base.ErrEventNotExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
