package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hospital-middleware-system/src/config"
	"hospital-middleware-system/src/database"
	"hospital-middleware-system/src/logger"
	"hospital-middleware-system/src/router"
)

// @title           Hospital Middleware System API
// @version         1.0
// @description     REST API middleware for multi-hospital deployments managing hospitals, staff, and patients.
// @host            localhost:8080
// @BasePath        /
// @schemes         http https
//
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT Bearer token. Format: "Bearer <token>"
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Init(config.AppConfig.AppEnv)

	if err := database.Connect(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	r := router.New(database.DB)

	addr := ":" + config.AppConfig.AppPort
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      1 * time.Minute,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Log.Info().
			Str("env", config.AppConfig.AppEnv).
			Str("addr", addr).
			Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server failed: %w", err)
	case <-quit:
		logger.Log.Info().Msg("Shutdown signal received, closing server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Log.Info().Msg("Server exited gracefully")
	return nil
}
