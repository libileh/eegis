package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/libileh/eegis/common/db"
	"github.com/libileh/eegis/posts/internal/infra/config"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	if err := LoadProperties(); err != nil {
		logger.Fatalf("Error loading environment file: %v", err)
	}

	if err := run(logger); err != nil {
		panic(err)
	}
}

func LoadProperties() error {
	// Check if the .env file exists
	localEnv := ".env"
	if _, err := os.Stat(localEnv); os.IsNotExist(err) {
		// if env file not found, load from home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error retrieving home directory: %v", err)
		}

		envFilePath := filepath.Join(homeDir, ".eegis.env")
		if err := godotenv.Load(envFilePath); err != nil {
			return fmt.Errorf("error loading .eegis.env file: %v", err)
		}
	} else {
		// Local load from .env file
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("error loading .env file: %v", err)
		}
	}
	return nil
}

func run(logger *zap.SugaredLogger) error {
	props, err := config.ProcessPostProperties(logger)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	postApp, err := setupApplication(props, logger)
	if err != nil {
		return fmt.Errorf("failed to setup application: %v", err)
	}

	return postApp.Start()
}

func setupApplication(props *config.PostProperties, logger *zap.SugaredLogger) (*PostApp, error) {

	dbCnx, err := setupDatabase(props.CommonProps.DBURL, logger)
	if err != nil {
		return nil, err
	}

	return NewPostApp(logger, props, dbCnx), nil
}

func setupDatabase(dbURL string, logger *zap.SugaredLogger) (*sql.DB, error) {
	// Use the common database setup
	dbCnx, err := db.SetupDatabase(dbURL, logger)
	if err != nil {
		strErr := fmt.Sprintf("failed to setup database connection: %v", err)
		logger.Fatalf(strErr)
		return nil, fmt.Errorf(strErr)
	}
	return dbCnx, nil
}
