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
	rows, err := r.pool.Query(ctx, `select * from event where event_id=$1`, eventID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, base.ErrEventNotExists)
	}
	defer rows.Close()

	events, err := converter.EventRowsToModel(rows)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("%s: %w", op, base.ErrEventNotExists)
	}

	return &events[0], nil
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

	sql := `select * from event`
	switch status {
	case "active":
		sql = `select * from event where canceled_at is null`
	case "completed":
		sql = `select * from event where canceled_at is not null`
	}

	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, base.ErrEventNotExists)
	}
	defer rows.Close()

	events, err := converter.EventRowsToModel(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, base.ErrEventNotExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
