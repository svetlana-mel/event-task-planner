package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/svetlana-mel/event-task-planner/internal/config"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres/sql"
	// base "github.com/svetlana-mel/event-task-planner/internal/repository"
)

// var _ base.PlannerRepository = (*Repository)(nil)

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(ctx context.Context, cfg config.DataBase) (*repository, error) {
	op := "repository.postgres.NewRepository"

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.Name,
	)

	dbpool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer dbpool.Close()

	// prepare and execute statements
	sqlFilesPath := []string{
		sql.CreateTablesStmt,
		sql.AddIndexesStmt,
	}

	for _, stmt := range sqlFilesPath {
		_, err = dbpool.Exec(ctx, stmt)
		if err != nil {
			return nil, fmt.Errorf("%s error execute stmt: %w", op, err)
		}
	}

	return &repository{dbpool}, nil
}
