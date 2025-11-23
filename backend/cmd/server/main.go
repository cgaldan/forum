package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum-backend/internal/api/router"
	"forum-backend/internal/config"
	"forum-backend/internal/repository"
	"forum-backend/internal/service"
	"forum-backend/internal/websocket"
	"forum-backend/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize logger
	appLogger := logger.NewLogger(os.Stdout, logger.InfoLevel)
	appLogger.Info("Starting Forum Application...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal("Failed to load configuration", "error", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		appLogger.Fatal("Invalid configuration", "error", err)
	}

	appLogger.Info("Configuration loaded successfully", "environment", cfg.Environment)

	// Initialize database
	db, err := repository.NewDatabase(cfg.Database.Path)
	if err != nil {
		appLogger.Fatal("Failed to initialize database", "error", err)
	}
	defer db.Close()

	appLogger.Info("Database initialized", "path", cfg.Database.Path)

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		appLogger.Fatal("Failed to run migrations", "error", err)
	}

	appLogger.Info("Database migrations completed")

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize WebSocket hub
	hub := websocket.NewHub(appLogger, repos.User)
	go hub.Run()

	appLogger.Info("WebSocket hub started")

	// Initialize services
	services := service.NewServices(repos, hub, appLogger)

	// Initialize router
	r := router.NewRouter(services, hub, cfg, appLogger)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", "port", cfg.Server.Port, "environment", cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server is shutting down...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exited gracefully")
}

