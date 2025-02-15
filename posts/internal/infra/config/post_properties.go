package config

import (
	"github.com/libileh/eegis/common/properties"
	"go.uber.org/zap"
	"time"
)

type PostProperties struct {
	Port            string `envconfig:"PORT" default:"9102" required:"true"`
	Version         string `envconfig:"VERSION" default:"0.1"`
	CommonProps     properties.CommonProperties
	PostsServiceURL string `envconfig:"POSTS_SERVICE_URL" default:"http://localhost:9102" required:"true"`
	UsersServiceURL string `envconfig:"USERS_SERVICE_URL" default:"http://localhost:9101" required:"true"`
}

func ProcessPostProperties(logger *zap.SugaredLogger) (*PostProperties, error) {
	props, err := properties.LoadProperties(logger, &PostProperties{})
	if err != nil {
		return nil, err
	}

	props.Port = ":" + props.Port
	props.CommonProps.AuthProps.Token.Exp *= time.Hour * 24
	return props, nil
}
