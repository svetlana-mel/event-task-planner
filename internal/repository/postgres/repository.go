package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/svetlana-mel/event-task-planner/internal/config"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres/migrations"
	// base "github.com/svetlana-mel/event-task-planner/internal/repository"
)

// var _ base.PlannerRepository = (*Repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) Close() {
	defer r.pool.Close()
}

func NewRepository(ctx context.Context, cfg config.DataBase) (*repository, error) {
	op := "repository.postgres.NewRepository"

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.Name,
	)

	fmt.Println(databaseUrl)

	dbpool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// prepare and execute statements
	sqlFilesPath := [][2]string{
		{"drop", migrations.DropTablesStmt},
		{"create", migrations.CreateTablesStmt},
		{"indexes", migrations.AddIndexesStmt},
	}

	for _, mig := range sqlFilesPath {
		info := mig[0]
		stmt := mig[1]
		_, err := dbpool.Exec(ctx, stmt)
		if err != nil {
			return nil, fmt.Errorf("%s error execute stmt %s: %w", op, info, err)
		}
	}

	return &repository{dbpool}, nil
}
