package main

import (
	"github.com/joho/godotenv"
	"github.com/libileh/eegis/notifications/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}
	if err := run(logger); err != nil {
		panic(err)
	}
}

func run(logger *zap.SugaredLogger) error {
	notifProps, err := config.ProcessNotificationProperties(logger)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}
	notifApp := NewNotificationApp(logger, notifProps)
	return notifApp.Start()
}
