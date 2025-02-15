package config

import (
	commonProps "github.com/libileh/eegis/common/properties"
	"go.uber.org/zap"
	"time"
)

type EmailProperties struct {
	FromEmail string        `envconfig:"MAILTRAP_FROM_EMAIL" default:"noreply@qoraalhub.com"`
	Exp       time.Duration `envconfig:"MAILTRAP_EXP" default:"10m"`
	Mailtrap  mailtrapProperties
}

type mailtrapProperties struct {
	ApiKey string `envconfig:"MAILTRAP_API_KEY" required:"true"`
	Url    string `envconfig:"MAILTRAP_URL" default:"https://api.mailtrap.io/v1/send"`
}

type NotificationProperties struct {
	CommonProps *commonProps.CommonProperties
	Port        string `envconfig:"NOTIF_PORT" default:"9103" required:"true"`
	Version     string `envconfig:"VERSION" default:"0.1"`
	Email       *EmailProperties
}

// ProcessEnvLoaderProperties processes environment variables into NotificationProperties.
func ProcessNotificationProperties(logger *zap.SugaredLogger) (*NotificationProperties, error) {
	notifProps, err := commonProps.LoadProperties(logger, &NotificationProperties{})
	if err != nil {
		return nil, err
	}

	notifProps.Port = ":" + notifProps.Port
	notifProps.Email.Exp *= time.Minute
	return notifProps, nil
}
