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
		application.Logger.Info("error run application", sl.AddErrorAtribute(err))
	}
	application.Logger.Info("server stopped")

	<-idleConnsClosed
}
