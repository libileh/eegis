package api

import (
	"expvar"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	common "github.com/libileh/eegis/common/app"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/monitoring"
	"github.com/libileh/eegis/common/ratelimiter"
	"github.com/libileh/eegis/topics/application"
	"github.com/libileh/eegis/topics/internal/infra/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type TopicApi struct {
	Properties    *config.TopicProperties
	Service       *application.Service
	Logger        *zap.SugaredLogger
	RateLimiter   *ratelimiter.RateLimiter
	Validator     *validator.Validate
	HttpError     *errors.Error
	Monitoring    *monitoring.Monitoring
	Authenticator *auth.JWTAuthenticator
}

func NewTopicApi(props *config.TopicProperties, service *application.Service, logger *zap.SugaredLogger, baseApp *common.App) *TopicApi {
	return &TopicApi{
		Properties:    props,
		Service:       service,
		Logger:        logger,
		RateLimiter:   baseApp.RateLimiter,
		Validator:     baseApp.Validate,
		HttpError:     baseApp.Error,
		Monitoring:    baseApp.Monitoring,
		Authenticator: baseApp.Authenticator,
	}

}

func (api *TopicApi) Mount() http.Handler {
	router := chi.NewRouter()

	// Middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	// Basic CORS
	router.Use(cors.Handler(cors.Options{
		//todo : adjust to frontend url
		AllowedOrigins: []string{api.Properties.CommonProps.FrontendUrl, "http://localhost"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(api.RateLimiter.RateLimiterMiddleware)
	router.With(api.Authenticator.BasicAuthMiddleware()).Get("/v1/metrics", expvar.Handler().ServeHTTP)

	api.mountTopicRoutes(router)
	return router
}

func (api *TopicApi) mountTopicRoutes(r chi.Router) {
	r.Route("/v1", func(r1 chi.Router) {
		r1.Use(api.Authenticator.AuthMiddleware)
		r1.Get("/users/{followerId}/topics", api.GetUserFollowedTopicsHandler) //todo: move to users service
		r1.Route("/topics", func(r2 chi.Router) {
			r2.Get("/", api.GetAllTopicsHandler)
			r2.Post("/", api.CreateTopicHandler)
			r2.Get("/{name}", api.GetTopicByUsernameHandler)
			r2.Route("/{topicId}", func(r3 chi.Router) {
				r3.Post("/followers/{userId}", api.FollowTopicHandler)
			})
		})
	})
}

func (api *TopicApi) HandlerError(w http.ResponseWriter, r *http.Request, customErr *errors.CustomError) {
	if customErr == nil {
		return
	}

	switch customErr.ErrType {
	case errors.NotFound:
		api.HttpError.NotFoundResponse(w, r, fmt.Sprintf("not found %v", customErr))
	case errors.BadRequest:
		api.HttpError.BadRequestResponse(w, r, customErr.Error())
	default:
		api.HttpError.InternalServerError(w, r, customErr.Error())
	}
}
