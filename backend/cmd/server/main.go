package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real-time-forum/internal/api/router"
	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/service"
	"real-time-forum/internal/websocket"
	"real-time-forum/packages/logger"
)

func main() {
	appLogger := logger.NewLogger(os.Stdout, logger.InfoLevel)
	appLogger.Info("Starting Real-Time Forum Application...")

	config, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatal("Failed to load configuration", "error", err)
	}

	if err := config.Validate(); err != nil {
		appLogger.Fatal("Invalid configuration", "error", err)
	}

	appLogger.Info("Configuration loaded successfully", "environment", config.Environment)

	db, err := database.NewDatabase(config.Database.Path)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	appLogger.Info("Database initialized", "path", config.Database.Path)

	if err := database.RunMigrations(db); err != nil {
		appLogger.Fatal("Failed to run database migrations", "error", err)
	}

	appLogger.Info("Database migrations completed successfully")

	repos := repository.NewRepositories(db)

	hub := websocket.NewHub(appLogger, repos.User)
	go hub.Run()

	appLogger.Info("WebSocket hub initialized")

	services := service.NewServices(repos, hub, appLogger)

	router := router.NewRouterNoGorilla(services, config, hub, appLogger)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	go func() {
		appLogger.Info("Starting server", "port", config.Server.Port, "environment", config.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Failed to shutdown server gracefully", "error", err)
	}

	appLogger.Info("Server stopped successfully")
}
