package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/svetlana-mel/event-task-planner/internal/app"
	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"
)

const (
	ENV_LOCAL = "local"
	ENV_DEV   = "dev"
	ENV_PROD  = "prod"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx, ENV_LOCAL)
	if err != nil {
		log.Fatalf("failed to init application: %s", err)
		os.Exit(1)
	}

	// Implementing Graceful Shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		// wait for interrupt signal
		<-sigint

		// We received an interrupt signal, close db pool and shut down
		application.Close()

		close(idleConnsClosed)
	}()

	err = application.Run()
	if err != nil {
		application.Logger.Error("error run application", sl.AddErrorAtribute(err))
	}
	application.Logger.Error("server stopped")

	<-idleConnsClosed

	// cfg := config.NewConfig(ENV_LOCAL)

	// log := setupLogger(ENV_LOCAL)

	// log.Info("starting event-task-planner", slog.String("env", cfg.Env))

	// storage, err := postgres.NewRepository(ctx, cfg.DataBase)
	// if err != nil {
	// 	log.Error("failed to init storage", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }
	// defer storage.Close()

	// events, err := storage.GetAllEvents(ctx, "all")
	// if err != nil {
	// 	log.Error("failed to get event", sl.AddErrorAtribute(err))
	// 	os.Exit(1)
	// }

	// fmt.Println(events)

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
