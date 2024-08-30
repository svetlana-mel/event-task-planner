package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/svetlana-mel/event-task-planner/internal/models"
	base "github.com/svetlana-mel/event-task-planner/internal/repository"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres/converter"
)

func (r *repository) CreateUser(
	ctx context.Context,
	username string,
	email string,
	passHash []byte,
) (uint64, error) {
	op := "repository.postgres.CreateUser"

	timestamp := time.Now()

	row := r.pool.QueryRow(ctx,
		`insert into "user" 
		(name, email, pass_hash, 
		created_date_time)
		values 
		(	$1, $2, $3,
			$4
		) RETURNING user_id`,
		username,
		email,
		passHash,
		timestamp,
	)

	var userID uint64
	err := row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "email_unique" {
			return 0, fmt.Errorf("%s: %w", op, base.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (r *repository) GetUser(
	ctx context.Context,
	email string,
) (*models.User, error) {
	op := "repository.postgres.GetUser"

	rows, _ := r.pool.Query(ctx, `select * from "user" where email=$1`, email)
	defer rows.Close()

	users, err := converter.UserRowsToModel(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, base.ErrUserNotExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("%s: %w", op, base.ErrUserNotExists)
	}

	return &users[0], nil
}
