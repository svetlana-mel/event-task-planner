// posgresql rows types
package converter

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
)

func TaskRowToModel(row pgx.Row) (*models.Task, error) {
	collectRow := row.(pgx.CollectableRow)
	task, err := pgx.RowToStructByName[models.Task](collectRow)
	if err != nil {
		return nil, fmt.Errorf("error in pgx.RowToStructByName: %w", err)
	}

	return &task, nil
}

func TaskRowsToModel(rows pgx.Rows) ([]models.Task, error) {
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Task])
	if err != nil {
		return nil, fmt.Errorf("error collect rows: %w", err)
	}

	return tasks, nil
}
