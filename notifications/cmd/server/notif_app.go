package main

import (
	"github.com/libileh/eegis/common/app"
	"github.com/libileh/eegis/notifications/application"
	"github.com/libileh/eegis/notifications/internal/config"
	"github.com/libileh/eegis/notifications/internal/infra/api"
	"github.com/libileh/eegis/notifications/internal/infra/sender"
	"github.com/libileh/eegis/notifications/internal/infra/template"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type NotificationApp struct {
	*app.App
	Properties *config.NotificationProperties
	Server     *http.Server
	Handler    http.Handler
	Service    *application.ServiceManager
}

// NewNotificationApp initializes a new NotificationApp
func NewNotificationApp(logger *zap.SugaredLogger, props *config.NotificationProperties) *NotificationApp {
	// Initialize the common App
	var baseApp = app.NewApp(logger, props.CommonProps, props.Version, nil)

	// Initialize email sender
	templManager := template.NewTemplateManager(logger)
	renderer := template.NewContentRenderer(logger)
	emailSender := sender.NewEmailSender(logger, props.Email)

	// Initialize notification service
	notificationService := application.NewNotificationService(emailSender, templManager, renderer, logger)

	// Initialize service manager
	serviceManager := application.NewServiceManager(notificationService)

	// router
	apiHandler := api.NotificationApi{
		Props:         props,
		ServerManager: serviceManager,
		Auth:          baseApp.Authenticator,
		RateLimiter:   baseApp.RateLimiter,
		Validate:      baseApp.Validate,
		HttpError:     baseApp.Error,
		Monitoring:    baseApp.Monitoring,
	}
	router := apiHandler.Mount()
	server := createServer(props, router)

	// Create the NotificationApp
	return &NotificationApp{
		App:        baseApp,
		Properties: props,
		Server:     server,
		Service:    serviceManager,
		Handler:    router,
	}
}

func createServer(props *config.NotificationProperties, router http.Handler) *http.Server {
	server := &http.Server{
		Addr:         props.Port,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	return server
}

func (n *NotificationApp) Start() error {
	n.Logger.Infow("NotificationApp server starting", "addr", n.Server.Addr)
	return n.Server.ListenAndServe()
}
