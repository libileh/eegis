package main

import (
	"database/sql"
	common "github.com/libileh/eegis/common/app"
	"github.com/libileh/eegis/topics/application"
	"github.com/libileh/eegis/topics/internal/client"
	"github.com/libileh/eegis/topics/internal/infra/api"
	"github.com/libileh/eegis/topics/internal/infra/config"
	"github.com/libileh/eegis/topics/internal/infra/persistence/repository"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type TopicApp struct {
	*common.App
	Server     *http.Server
	Handler    http.Handler         // Chi router for post-specific routes
	Service    *application.Service // Post service for database operations
	Properties *config.TopicProperties
}

func NewTopicApp(logger *zap.SugaredLogger, props *config.TopicProperties, db *sql.DB) *TopicApp {
	baseApp := common.NewApp(logger, &props.CommonProps, props.Version, db)
	store := repository.NewPostgresStorage(db)
	userClient := client.NewUserClient(props.UserBaseURL, props)
	service := application.NewService(store.Topics(), userClient)
	apiHandler := api.NewTopicApi(props, service, logger, baseApp)
	router := apiHandler.Mount()
	server := &http.Server{
		Addr:         props.Port,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &TopicApp{
		App:        baseApp,
		Server:     server,
		Handler:    router,
		Service:    service,
		Properties: props,
	}
}

func (t *TopicApp) Start() error {
	t.Logger.Info("TopicApp server starting... ", "addr", t.Server.Addr)
	return t.Server.ListenAndServe()
}
