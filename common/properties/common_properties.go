package properties

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"time"
)

type TokenProperties struct {
	Secret string        `envconfig:"AUTH_TOKEN_SECRET" default:"secret"`
	Exp    time.Duration `envconfig:"AUTH_TOKEN_EXP" default:"72h"`
}

type AuthProperties struct {
	Username string `envconfig:"AUTH_USERNAME" default:"admin"`
	Password string `envconfig:"AUTH_PASSWORD" default:"password"`
	Audience string `envconfig:"AUTH_AUDIENCE" default:"qoraalhub"`
	Issuer   string `envconfig:"AUTH_ISSUER" default:"qoraalhub"`
	Token    TokenProperties
}

type RateLimiterProperties struct {
	RequestsPerSecond int           `envconfig:"REQUEST_PER_SECOND" default:"5"`
	TimeFrequency     time.Duration `envconfig:"TIME_FREQUENCY" default:"5s"`
	Enabled           bool          `envconfig:"RATE_LIMITER_ENABLED" default:"true"`
}

type CommonProperties struct {
	DBURL       string `envconfig:"DB_URL"` //todo add ENABLED_DB
	Env         string `envconfig:"ENV" default:"DEV"`
	FrontendUrl string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`
	AuthProps   AuthProperties
	RateLimiter RateLimiterProperties
	TZ          string `envconfig:"TZ" default:"Europe/Paris"`
}

// LoadProperties processes environment variables into the given target struct.
// It can be used for any struct following envconfig conventions.
func LoadProperties[T any](logger *zap.SugaredLogger, target *T) (*T, error) {
	if err := envconfig.Process("", target); err != nil {
		logger.Error(fmt.Sprintf("error loading properties: %v", err))
		return nil, fmt.Errorf("failed to load properties for %T: %w", target, err)
	}

	logger.Infof("Properties loaded successfully for %T", target)
	return target, nil
}
