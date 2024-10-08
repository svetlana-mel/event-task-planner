package converter

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
)

func TaskRowsToModel(rows pgx.Rows) ([]models.Task, error) {
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Task])
	if err != nil {
		return nil, fmt.Errorf("error collect rows: %w", err)
	}

	return tasks, nil
}
