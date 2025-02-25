package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/libileh/eegis/notifications/internal/config"
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
		//local load from .env file
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("error loading .env file: %v", err)
		}
	}
	return nil
}

func run(logger *zap.SugaredLogger) error {
	notifProps, err := config.ProcessNotificationProperties(logger)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}
	notifApp := NewNotificationApp(logger, notifProps)
	return notifApp.Start()
}
