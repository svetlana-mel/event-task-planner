package app

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/svetlana-mel/event-task-planner/internal/config"
	"github.com/svetlana-mel/event-task-planner/internal/lib/jwt"
	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"

	"github.com/svetlana-mel/event-task-planner/internal/repository"
	"github.com/svetlana-mel/event-task-planner/internal/repository/postgres"

	"github.com/svetlana-mel/event-task-planner/internal/server/handlers/auth"
	"github.com/svetlana-mel/event-task-planner/internal/server/handlers/event"
	"github.com/svetlana-mel/event-task-planner/internal/server/handlers/task"
	"github.com/svetlana-mel/event-task-planner/internal/server/router"

	auth_middleware "github.com/svetlana-mel/event-task-planner/internal/server/middleware"
	auth_service "github.com/svetlana-mel/event-task-planner/internal/services/auth"
)

type Keys struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

type App struct {
	Env string

	Config *config.Config

	JWTKeys Keys

	Logger *slog.Logger

	repository repository.PlannerRepository

	HttpServer *http.Server

	closeFuncs []func()
}

func NewApp(ctx context.Context, env string) (*App, error) {
	app := &App{Env: env}

	app.closeFuncs = make([]func(), 0)

	err := app.initDependencies(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	a.Logger.Info("run http server")
	return a.HttpServer.ListenAndServe()
}

func (a *App) initDependencies(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initJwtKeys,
		a.initRepository,
		a.initHttpServer,
	}

	for _, initFn := range inits {
		err := initFn(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	cfg := config.NewConfig(a.Env)
	a.Config = cfg
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	log := sl.SetupLogger(a.Config.Env)
	a.Logger = log
	return nil
}

func (a *App) initJwtKeys(_ context.Context) error {
	publicKey, err := jwt.LoadPublicKey(a.Config.PublicKeyPath)
	if err != nil {
		a.Logger.Error("failed to init JWTKeys: %w", sl.AddErrorAtribute(err))
		return err
	}

	privateKey, err := jwt.LoadPrivateKey(a.Config.PrivateKeyPath)
	if err != nil {
		a.Logger.Error("failed to init JWTKeys: %w", sl.AddErrorAtribute(err))
		return err
	}

	a.JWTKeys.Public = publicKey
	a.JWTKeys.Private = privateKey

	return nil
}

func (a *App) initRepository(ctx context.Context) error {
	rep, err := postgres.NewRepository(ctx, a.Config.DataBase)
	if err != nil {
		a.Logger.Error("failed to init storage", sl.AddErrorAtribute(err))
		return fmt.Errorf("error init storage: %w", err)
	}

	a.closeFuncs = append(a.closeFuncs, func() {
		rep.Close()
	})
	a.Logger.Info("success init repository")
	a.repository = rep
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.URLFormat)

	cfg := a.Config.HTTPServer

	authService := auth_service.New(
		a.JWTKeys.Private,
		a.Logger,
		a.repository,
		a.repository,
		a.Config.JwtTTL,
	)

	authHandler := &auth.Handler{
		Auth:   authService,
		Logger: a.Logger,
	}

	mux.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/signup", authHandler.Signup)
	})

	eventHandler := &event.Handler{
		Repo:   a.repository,
		Logger: a.Logger,
	}

	taskHandler := &task.Handler{
		Repo:   a.repository,
		Logger: a.Logger,
	}

	mux.Group(func(r chi.Router) {
		// middleware that verifies the token used only for product endpoints
		r.Use(auth_middleware.New(a.JWTKeys.Public, a.Logger))

		router.SetupRoutes(r, eventHandler, taskHandler, authHandler)
	})

	srv := http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	a.HttpServer = &srv

	// prepare server shutdown closer func
	a.closeFuncs = append(a.closeFuncs, func() {
		if err := a.HttpServer.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
	})

	a.Logger.Info("success setup http server")

	return nil
}

func (a *App) Close() {
	for _, close := range a.closeFuncs {
		close()
	}
	a.Logger.Info("closed db connection pool and http server")
}
