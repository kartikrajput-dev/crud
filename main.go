package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kartikrajput-dev/crud/internal/handler"
	_ "github.com/lib/pq"

	"github.com/KARTIKrocks/apikit/dbx"
	"github.com/KARTIKrocks/apikit/health"
	"github.com/KARTIKrocks/apikit/middleware"
	"github.com/KARTIKrocks/apikit/router"
	"github.com/KARTIKrocks/apikit/server"
	appcfg "github.com/kartikrajput-dev/crud/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := appcfg.Load()
	if err != nil {
		logger.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		logger.Error("failed to open db", "err", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	dbx.SetDefault(db)

	// Health checker
	h := health.NewChecker(health.WithTimeout(3 * time.Second))
	h.AddCheck("postgres", func(ctx context.Context) error {
		return db.PingContext(ctx)
	})

	// Router
	r := router.New()
	r.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recover(),
		middleware.SecureHeaders(),
		middleware.Timeout(30*time.Second),
	)

	r.Get("/health", h.Handler())
	r.Get("/health/live", h.LiveHandler())

	api := r.Group("/api/v1")
	api.Get("/users", handler.ListUsers)
	api.Post("/users", handler.CreateUser)
	api.Get("/users/{id}", handler.GetUser)
	api.Put("/users/{id}", handler.UpdateUser)
	api.Delete("/users/{id}", handler.DeleteUser)

	srv := server.New(r, server.WithAddr(cfg.ServerAddr), server.WithLogger(logger))

	srv.OnStart(func() error {
		if err := db.PingContext(context.Background()); err != nil {
			return fmt.Errorf("postgres not reachable: %w", err)
		}
		logger.Info("connected to postgres")
		return nil
	})

	srv.OnShutdown(func(ctx context.Context) error {
		logger.Info("closing database connection")
		return db.Close()
	})

	logger.Info("starting server", "addr", cfg.ServerAddr)
	if err := srv.Start(); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
