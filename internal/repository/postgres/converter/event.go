package converter

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
)

func EventRowsToModel(rows pgx.Rows) ([]models.Event, error) {
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Event])
	if err != nil {
		return nil, fmt.Errorf("error collect rows: %w", err)
	}

	return events, nil
}
