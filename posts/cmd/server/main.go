package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/libileh/eegis/common/db"
	"github.com/libileh/eegis/posts/internal/infra/config"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//load properties from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}
	if err := run(logger); err != nil {
		log.Fatal(err)
	}
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
