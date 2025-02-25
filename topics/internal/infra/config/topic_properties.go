package config

import (
	commonProps "github.com/libileh/eegis/common/properties"
	"go.uber.org/zap"
)

type TopicProperties struct {
	CommonProps commonProps.CommonProperties
	Port        string `envconfig:"PORT" default:"9104" required:"true"`
	Version     string `envconfig:"VERSION" default:"0.1"`
	UserBaseURL string `envconfig:"USER_BASE_URL" required:"true"`
	AuthProps   AuthClientProperties
}

type AuthClientProperties struct {
	AuthToken    string `envconfig:"AUTH_TOKEN" default:"" required:"true"`
	AuthEndpoint string `envconfig:"AUTH_ENDPOINT" default:"/v1/auth/token/refresh" required:"true"`
}

func ProcessTopicsProperties(logger *zap.SugaredLogger) (*TopicProperties, error) {
	topicsProps, err := commonProps.LoadProperties(logger, &TopicProperties{})
	if err != nil {
		return nil, err
	}

	topicsProps.Port = ":" + topicsProps.Port
	topicsProps.AuthProps.AuthEndpoint = topicsProps.UserBaseURL + topicsProps.AuthProps.AuthEndpoint
	return topicsProps, nil

}
