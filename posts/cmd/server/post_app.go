package main

import (
	"database/sql"
	common "github.com/libileh/eegis/common/app"
	"github.com/libileh/eegis/posts/application"
	"github.com/libileh/eegis/posts/internal/infra/api"
	"github.com/libileh/eegis/posts/internal/infra/config"
	"github.com/libileh/eegis/posts/internal/infra/persistence/repository"
	"github.com/libileh/eegis/posts/pkg/client"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// PostApp extends the common App and adds post-specific components.
type PostApp struct {
	*common.App // Embed the common App
	Server      *http.Server
	Handler     http.Handler                // Chi router for post-specific routes
	Service     *application.ServiceManager // Post service for database operations
	Properties  *config.PostProperties
}

// NewPostApp initializes a new PostApp with user-specific components.
func NewPostApp(logger *zap.SugaredLogger, props *config.PostProperties, db *sql.DB) *PostApp {
	// Initialize the common App
	baseApp := common.NewApp(logger, &props.CommonProps, props.Version, db)

	//Initialize Repositories DB
	store := repository.NewPostgresStorage(db)

	//Initialize service
	postService := application.NewPostService(store.Posts(), store.Comments(), store.Feeds(), baseApp.EventBus)
	userService := client.NewHttpUserService(props.UsersServiceURL)

	serviceMgmt := application.NewServiceManager(postService, userService)

	//Initialize router
	apiHandler := api.NewPostApi(props, serviceMgmt, logger, baseApp)
	router := apiHandler.Mount()

	server := createServer(props, router)
	// Create the UserApp
	return &PostApp{
		App:        baseApp,
		Server:     server,
		Handler:    router,
		Service:    serviceMgmt,
		Properties: props,
	}
}

func createServer(props *config.PostProperties, router http.Handler) *http.Server {
	server := &http.Server{
		Addr:         props.Port,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	return server
}

// Start initializes the UserApp and starts the server.
func (p *PostApp) Start() error {
	// Start the server
	p.Logger.Infow("PostApp server starting", "addr", p.Server.Addr)
	return p.Server.ListenAndServe()
}
