package converter

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
	base "github.com/svetlana-mel/event-task-planner/internal/repository"
)

func EventRowToModel(row pgx.CollectableRow) (models.Event, error) {
	var event models.Event
	var tasksJSON []byte

	err := row.Scan(
		&event.EventID,
		&event.Name,
		&event.Description,
		&event.DateTime,
		&event.CanceledAt,
		&event.DeletedAt,
		&event.FkUserID,
		&tasksJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Event{}, base.ErrEventNotExists
		}
		return models.Event{}, err
	}

	err = json.Unmarshal(tasksJSON, &event.Tasks)
	if err != nil {
		return models.Event{}, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return event, nil
}

func EventsRowsToModel(rows pgx.Rows) ([]models.Event, error) {
	events, err := pgx.CollectRows(rows, EventRowToModel)
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return nil, err
	}

	return events, nil
}
