package app

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/monitoring"
	"github.com/libileh/eegis/common/properties"
	"github.com/libileh/eegis/common/ratelimiter"
	"go.uber.org/zap"
)

type App struct {
	Logger        *zap.SugaredLogger
	Properties    *properties.CommonProperties
	DB            *sql.DB
	Error         *errors.Error
	L10nUtils     *L10nUtils
	Validate      *validator.Validate
	Authenticator *auth.JWTAuthenticator
	RateLimiter   *ratelimiter.RateLimiter
	Monitoring    *monitoring.Monitoring
}

func NewApp(logger *zap.SugaredLogger, loadedProperties *properties.CommonProperties, version string, db *sql.DB) *App {

	err := &errors.Error{
		Logger: logger,
	}

	l10n := &L10nUtils{
		Logger:     logger,
		Properties: loadedProperties,
	}
	if err := l10n.SetTimeZone(); err != nil {
		logger.Warnf("%s", err)
	}

	//validator
	validate := validator.New(validator.WithRequiredStructEnabled())

	//auth
	authenticator := auth.NewJWTAuthenticator(&loadedProperties.AuthProps, err, validate)

	//rate Limiter
	ratelimit := &ratelimiter.RateLimiter{
		Properties:      &loadedProperties.RateLimiter,
		RateLimiterAlgo: ratelimiter.NewFixedWindowLimiter(loadedProperties.RateLimiter.RequestsPerSecond, loadedProperties.RateLimiter.TimeFrequency),
		Error:           err,
	}

	//monitoring
	monitor := &monitoring.Monitoring{
		Metrics: monitoring.NewMetrics(version, db),
		Health:  monitoring.NewHealth(version, loadedProperties.Env),
	}

	return &App{
		Logger:        logger,
		Properties:    loadedProperties,
		DB:            db,
		Error:         err,
		L10nUtils:     l10n,
		Validate:      validate,
		Authenticator: authenticator,
		RateLimiter:   ratelimit,
		Monitoring:    monitor,
	}
}
