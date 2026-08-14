package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"hospital-middleware-system/src/config"
	"hospital-middleware-system/src/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStr := config.AppConfig.DBConnectionString()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	DB = db
	logger.Log.Info().Str("db", config.AppConfig.DBName).Msg("Database connected successfully")
	return nil
}

func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Error closing database")
		} else {
			logger.Log.Info().Msg("Database connection closed")
		}
	}
}
