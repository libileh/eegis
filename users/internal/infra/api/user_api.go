package api

import (
	"expvar"
	"fmt"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/monitoring"
	"github.com/libileh/eegis/common/ratelimiter"
	"github.com/libileh/eegis/users/application"
	"github.com/libileh/eegis/users/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UserApi struct {
	Properties  *config.UserProperties
	Service     *application.ServiceManager
	Logger      *zap.SugaredLogger
	Auth        *auth.JWTAuthenticator
	RateLimiter *ratelimiter.RateLimiter
	Validate    *validator.Validate
	HttpError   *errors.Error
	Monitoring  *monitoring.Monitoring
}

func (api *UserApi) Mount() http.Handler {
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
	router.Get("/health", api.Monitoring.Health.HealthCheckHandler)
	router.With(api.Auth.BasicAuthMiddleware()).Get("/v1/metrics", expvar.Handler().ServeHTTP)

	// Mount specific API routes
	api.mountUserRoutes(router)
	api.mountAuthRoutes(router)
	return router
}

// timerSleep tested for gracefully shutdown
func (api *UserApi) timerSleep(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		next.ServeHTTP(w, r)
	})
}

func (api *UserApi) mountUserRoutes(r chi.Router) {

	r.Route("/v1/users", func(ru chi.Router) {
		ru.Use(api.Auth.AuthMiddleware)
		ru.Get("/role-precedence", api.checkRolePrecedenceHandler) // New endpoint
		ru.Route("/{userId}", func(r chi.Router) {
			r.Get("/", api.GetUserByIdHandler)
			r.Post("/follow/{topicId}", api.followTopicHandler)
			r.Delete("/unfollow", api.unfollowTopicHandler)
			//r.Get("/feeds", api.getUserFeedHandler)
		})
		ru.Put("/activate/{token}", api.activateUserHandler)
	})
}

func (api *UserApi) mountAuthRoutes(r chi.Router) {
	r.Route("/v1/auth", func(ra chi.Router) {
		ra.Post("/users", api.RegisterUserHandler)
		ra.Post("/token", api.generateTokenHandler)
		ra.Get("/token/refresh", api.RefreshTokenHandler)
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
