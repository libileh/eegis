package config

import (
	commonProps "github.com/libileh/eegis/common/properties"
	"go.uber.org/zap"
	"time"
)

type TopicProperties struct {
	BaseUrl string `envconfig:"TOPIC_BASE_URL" default:"http://localhost:9104"`
}

type RedisProperties struct {
	Addr    string `envconfig:"REDIS_CACHE_ADDR" default:"gurilab.local:6379"`
	Pwd     string `envconfig:"REDIS_CACHE_PASSWORD"`
	DB      int    `envconfig:"REDIS_CACHE_DB" default:"0"`
	Enabled bool   `envconfig:"REDIS_CACHE_ENABLED" default:"false"`
}

type UserProperties struct {
	CommonProps       commonProps.CommonProperties
	Port              string `envconfig:"PORT" default:"9101" required:"true"`
	Version           string `envconfig:"VERSION" default:"0.1"`
	Redis             RedisProperties
	NotificationProps NotificationProps
	Topics            TopicProperties
}

type NotificationProps struct {
	FromEmail           string `envconfig:"MAILTRAP_FROM_EMAIL" default:"noreply@qoraalhub.com"`
	ApiKey              string `envconfig:"MAILTRAP_API_KEY" required:"true"`
	NotificationBaseUrl string `envconfig:"NOTIFICATION_BASE_URL" required:"true"`
}

// ProcessUserProperties processes environment variables into UserProperties and adjusts specific fields.
func ProcessUserProperties(logger *zap.SugaredLogger) (*UserProperties, error) {
	userProps, err := commonProps.LoadProperties(logger, &UserProperties{})
	if err != nil {
		return nil, err
	}

	// Apply custom field adjustments for CommonProps properties.
	userProps.Port = ":" + userProps.Port
	userProps.CommonProps.AuthProps.Token.Exp *= time.Hour * 24
	return userProps, nil
}
