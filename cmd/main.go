package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/svetlana-mel/event-task-planner/internal/config"
	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"
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

	event, err := storage.GetEvent(ctx, 1)
	if err != nil {
		log.Error("failed to get event", sl.AddErrorAtribute(err))
		os.Exit(1)
	}

	fmt.Println(event.Tasks)
	fmt.Println(event.Tasks[0])

	// log.Info("user created")

	// err = storage.CreateTask(ctx, &models.Task{
	// 	Name:     "first task",
	// 	FkUserID: userID,
	// })
	// if err != nil {
	// 	log.Error("failed to create task", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }

	// log.Info("task created")

	// tasks, err := storage.GetAllTasks(ctx, "any")
	// if err != nil {
	// 	log.Error("failed to get all tasks", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }

	// log.Info("got all tasks")

	// task := tasks[0]
	// id := task.TaskID

	// fmt.Println(task)

	// task.Description = "new description"

	// err = storage.UpdateTask(ctx, &task)
	// if err != nil {
	// 	log.Error("failed to update task", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }

	// log.Info("task updated")

	// t1, err := storage.GetTask(ctx, id)
	// if err != nil {
	// 	log.Error("failed to get task", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }
	// fmt.Println(t1)

	// log.Info("got task")

	// err = storage.SetTaskCompletionStatus(ctx, task.TaskID, true)
	// if err != nil {
	// 	log.Error("failed to set completed", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }
	// log.Info("completed set")

	// t2, err := storage.GetTask(ctx, id)
	// if err != nil {
	// 	log.Error("failed to get task", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }
	// fmt.Println(t2)
	// log.Info("got task")

	// err = storage.DeleteTask(ctx, task.TaskID)
	// if err != nil {
	// 	log.Error("failed to delete", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }
	// log.Info("deleted task")

	// t3, err := storage.GetTask(ctx, id)
	// if err != nil {
	// 	log.Error("failed to get task", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }
	// log.Info("got task")
	// fmt.Println(t3)
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
