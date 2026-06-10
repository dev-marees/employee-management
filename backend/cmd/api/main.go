// Command api is the HTTP entrypoint for the Employee Management System.
//
// @title                       Employee Management System API
// @version                     1.0
// @description                 Backend service for the EMS frontend (auth, employees, dashboard).
// @termsOfService              http://swagger.io/terms/
// @contact.name                API Support
// @contact.email               support@example.com
// @license.name                MIT
// @host                        localhost:8080
// @BasePath                    /api/v1
// @schemes                     http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT access token.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/ems/internal/app"
	"github.com/example/ems/internal/config"
	"github.com/example/ems/internal/database"
	"github.com/example/ems/internal/logger"
	"github.com/example/ems/internal/router"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		// Logger may not be initialized yet; fall back to stderr.
		panic(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log, err := logger.New(cfg.Log.Level, cfg.IsProduction())
	if err != nil {
		return err
	}
	defer func() { _ = log.Sync() }()

	log.Info("starting service",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	db, err := database.New(cfg.Database, log, cfg.IsProduction())
	if err != nil {
		return err
	}
	defer func() { _ = database.Close(db) }()

	if err := database.Migrate(db); err != nil {
		return err
	}
	log.Info("database migrated")

	container := app.New(cfg, db, log)
	engine := router.Setup(container)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start the server in a goroutine so shutdown signals can be handled.
	serverErr := make(chan error, 1)
	go func() {
		log.Info("http server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Wait for interrupt or a fatal server error.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		log.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	// Graceful shutdown with a bounded timeout.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
		return err
	}
	log.Info("server stopped cleanly")
	return nil
}
