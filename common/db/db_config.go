package db

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBConfig holds the configuration for the database connection.
type DBConfig struct {
	URL string
}

// SetupDatabase creates a new database connection.
func SetupDatabase(dbURL string, logger *zap.SugaredLogger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ping the database to verify the connection.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to the database")
	return db, nil
}
