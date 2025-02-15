package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/libileh/eegis/common/db"
	"github.com/libileh/eegis/users/internal/config"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//load properties from .env file
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	if err := run(logger); err != nil {
		log.Fatal(err)
	}
}

func run(logger *zap.SugaredLogger) error {
	userProps, err := config.ProcessUserProperties(logger)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}
	userApp, err := setupApplication(userProps, logger)
	if err != nil {
		return fmt.Errorf("failed to setup application: %v", err)
	}
	return userApp.Start()
}

func setupApplication(userProperties *config.UserProperties, logger *zap.SugaredLogger) (*UserApp, error) {

	dbCnx, err := setupDatabase(userProperties.CommonProps.DBURL, logger)
	if err != nil {
		return nil, err
	}

	return NewUserApp(logger, userProperties, dbCnx), nil
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
