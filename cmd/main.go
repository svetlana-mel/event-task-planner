package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/svetlana-mel/event-task-planner/internal/config"
	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"
	"github.com/svetlana-mel/event-task-planner/internal/models"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres"
)

const (
	ENV_LOCAL = "local"
	ENV_DEV   = "dev"
	ENV_PROD  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig(ENV_LOCAL)

	log := setupLogger(ENV_LOCAL)

	log.Info("starting event-task-planner", slog.String("env", cfg.Env))

	storage, err := postgres.NewRepository(ctx, cfg.DataBase)
	if err != nil {
		log.Error("failed to init storage", sl.AddErrorAtribute(err))
		os.Exit(1)
	}
	defer storage.Close()

	userID, err := storage.CreateTmpUser(ctx)
	if err != nil {
		log.Error("failed to create user", sl.AddErrorAtribute(err))
		os.Exit(1)
	}

	err = storage.CreateTask(ctx, models.Task{
		Name:     "first task",
		FkUserID: userID,
	})
	if err != nil {
		log.Error("failed to create task", sl.AddErrorAtribute(err))
		os.Exit(1)
	}

	tasks, err := storage.GetAllTasks(ctx, "any")
	if err != nil {
		log.Error("failed to get all task", sl.AddErrorAtribute(err))
		os.Exit(1)
	}

	fmt.Println(tasks)
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case ENV_LOCAL:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_DEV:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_PROD:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
