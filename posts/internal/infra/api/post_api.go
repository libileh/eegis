package api

import (
	"expvar"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	common "github.com/libileh/eegis/common/app"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/monitoring"
	"github.com/libileh/eegis/common/ratelimiter"
	"github.com/libileh/eegis/posts/application"
	"github.com/libileh/eegis/posts/internal/infra/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type PostApi struct {
	Properties    *config.PostProperties
	Service       *application.ServiceManager
	Logger        *zap.SugaredLogger
	RateLimiter   *ratelimiter.RateLimiter
	Validate      *validator.Validate
	HttpError     *errors.Error
	Monitoring    *monitoring.Monitoring
	Authenticator *auth.JWTAuthenticator
}

func NewPostApi(props *config.PostProperties, service *application.ServiceManager, logger *zap.SugaredLogger, baseApp *common.App) *PostApi {
	return &PostApi{
		Properties:    props,
		Service:       service,
		Logger:        logger,
		RateLimiter:   baseApp.RateLimiter,
		Validate:      baseApp.Validate,
		HttpError:     baseApp.Error,
		Monitoring:    baseApp.Monitoring,
		Authenticator: baseApp.Authenticator,
	}
}

func (api *PostApi) Mount() http.Handler {
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

	// Mount specific API routes
	api.mountPostRoutes(router)
	return router
}

func (api *PostApi) mountPostRoutes(r chi.Router) {

	r.Route("/v1/posts", func(rp chi.Router) {
		rp.Use(api.Authenticator.AuthMiddleware)
		rp.Post("/", api.CreatePostHandler)
		rp.Get("/", api.getAllPosts)

		rp.Route("/{postId}", func(r chi.Router) {
			r.Use(api.postLoaderMiddleware)
			r.Get("/", api.GetPostByIdHandler)
			r.Put("/", api.checkPostOwnership("moderator", api.updatePostHandler))
			r.Put("/review", api.checkPostOwnership("moderator", api.PostModerationHandler))
			r.Delete("/", api.checkPostOwnership("admin", api.deletePost))
		})
	})
}

// timerSleep tested for gracefully shutdown
func (api *PostApi) timerSleep(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		next.ServeHTTP(w, r)
	})
}

func ValidateID(r *http.Request, idParam string) (*uuid.UUID, error) {
	uuidParam := chi.URLParam(r, idParam)
	id, err := uuid.Parse(uuidParam)
	if err != nil {
		return nil, fmt.Errorf("invalid %s param: %v", idParam, err)
	}
	return &id, nil
}
