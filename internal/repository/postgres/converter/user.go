package converter

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/svetlana-mel/event-task-planner/internal/models"
)

func UserRowsToModel(rows pgx.Rows) ([]models.User, error) {
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("error collect rows: %w", err)
	}

	return users, nil
}
