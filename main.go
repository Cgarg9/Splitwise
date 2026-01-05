package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"splitwise-clone/db"
	"splitwise-clone/internal/app"
	"splitwise-clone/internal/logger"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

//	@title			Splitwise Clone API
//	@version		1.0
//	@description	This is a sample server for a Splitwise clone application.
//	@termsOfService	http://swagger.io/terms/

// host localhost:8080
//	@BasePath	/api/v1

// securityDefinitions.apikey BearerAuth
//
//	@in		header
//	@name	Authorization
func main() {
	// Setup logger with file output and rotation
	logConfig := logger.DefaultConfig()
	logConfig.Level = "debug" // Use "info" for production
	if err := logger.Setup(logConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	log.Info().Msg("Application started")

	ctx := context.Background()
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database pool")
	}
	defer pool.Close()

	log.Info().Msg("Database pool created successfully")

	// Create and start HTTP server
	httpServer := app.NewHTTPServer(8080, pool)

	// Run server in a goroutine
	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	log.Info().Msg("Server is running on http://localhost:8080")

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}
