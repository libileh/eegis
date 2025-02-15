package main

import (
	"database/sql"
	"github.com/libileh/eegis/common/app"
	"github.com/libileh/eegis/users/application"
	"github.com/libileh/eegis/users/internal/config"
	"github.com/libileh/eegis/users/internal/infra/api"
	"github.com/libileh/eegis/users/internal/infra/client"
	"github.com/libileh/eegis/users/internal/infra/persistence/cache"
	"github.com/libileh/eegis/users/internal/infra/persistence/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// UserApp extends the common App and adds user-specific components.
type UserApp struct {
	*app.App   // Embed the common App
	Properties *config.UserProperties
	Server     *http.Server
	Handler    http.Handler
	Service    *application.ServiceManager // User service Manager for DB/HTTP operations
	Cache      *redis.Client               // Redis cache for user data
}

// UserRepository handles database operations for the Users domain.

// NewUserApp initializes a new UserApp with user-specific components.
func NewUserApp(logger *zap.SugaredLogger, userProps *config.UserProperties, db *sql.DB) *UserApp {

	// Initialize the common App
	baseApp := app.NewApp(logger, &userProps.CommonProps, userProps.Version, db)

	//Repository DB
	store := repository.NewPostgresStorage(db)
	userService := application.NewUserService(store.Users(), store.Roles(), store.Followers())

	// password service
	passwordService := application.NewPasswordService()

	// init Redis cache
	var redisCache *redis.Client
	if userProps.Redis.Enabled {
		redisCache = cache.NewRedisClient(userProps.Redis.Addr, userProps.Redis.Pwd, userProps.Redis.DB)
		logger.Info("Redis Cache Enabled")
	}
	//Cache Repository
	redisStore := cache.NewRedisStorage(redisCache)
	userCache := application.NewUserCacheService(redisStore.Users())
	notifService := client.NewHttpNotifService(userProps)
	serviceMgmt := application.NewServiceManager(userService, userCache, passwordService, notifService)

	// Initialize UserApi with its dependencies
	apiHandler := &api.UserApi{
		Properties:  userProps,
		Service:     serviceMgmt,
		Logger:      logger,
		Auth:        baseApp.Authenticator,
		RateLimiter: baseApp.RateLimiter,
		Validate:    baseApp.Validate,
		HttpError:   baseApp.Error,
		Monitoring:  baseApp.Monitoring,
	}
	//router
	router := apiHandler.Mount()
	server := &http.Server{
		Addr:         userProps.Port,
		Handler:      router,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	// Create the UserApp
	return &UserApp{
		Properties: userProps,
		Server:     server,
		Handler:    router,
		App:        baseApp,
		Service:    serviceMgmt,
		Cache:      redisCache,
	}
}

// Start initializes the UserApp and starts the server.
func (ua *UserApp) Start() error {
	// Start the server
	ua.Logger.Infow("UserApp server starting", "addr", ua.Server.Addr)
	return ua.Server.ListenAndServe()
}
