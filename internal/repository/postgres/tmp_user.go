package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/svetlana-mel/event-task-planner/internal/models"
)

// CreateTmpUser: пока нет сервиса авторизации,
// создаем временного пользователя
// получаем его id и далее с id тестируем разрабатываемый функционал
func (r *repository) CreateTmpUser(ctx context.Context) (uint64, error) {
	op := "repository.postgres.CreateTmpUser"

	user := models.User{
		Name:     "tmp user",
		Email:    "some.email@mail.com",
		PassHash: []byte("my=passw"),
	}

	timeNow := time.Now().Local().UTC()
	fmt.Println(timeNow)

	_, err := r.pool.Exec(ctx,
		`insert into "user" 
		(name, email, pass_hash, created_date_time)
		values 
		($1, $2, $3, $4)`,
		user.Name, user.Email, user.PassHash, timeNow,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// get id
	row := r.pool.QueryRow(ctx, `select user_id from "user"`)
	var userID uint64
	row.Scan(&userID)

	return userID, nil
}
