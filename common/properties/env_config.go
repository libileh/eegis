package properties

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
	"os"
)

// LoadEnv loads environment variables from a .env file (if present)
// and into the provided list of structs (e.g., CommonProperties, PostProperties, etc.).
func LoadEnv(logger *zap.SugaredLogger, targets ...interface{}) error {
	// Load .env file if it exists.
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	// Use envconfig to populate each provided struct.
	for _, target := range targets {
		if err := envconfig.Process("", target); err != nil {
			logger.Errorf("Failed to process environment variables for target %T: %v", target, err)
			return fmt.Errorf("failed to load environment variables for target %T: %w", target, err)
		}
		logger.Infof("Loaded environment variables for target %T successfully", target)
	}

	log.Println("Environment variables loaded successfully")
	return nil
}

// GetEnv retrieves an environment variable or returns a default value.
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
