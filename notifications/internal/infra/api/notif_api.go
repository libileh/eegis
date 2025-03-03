package api

import (
	"expvar"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/monitoring"
	"github.com/libileh/eegis/common/ratelimiter"
	"github.com/libileh/eegis/notifications/application"
	"github.com/libileh/eegis/notifications/internal/config"
	"net/http"
	"time"
)

// NotificationApi handles HTTP requests for notifications
type NotificationApi struct {
	Props         *config.NotificationProperties
	ServerManager *application.ServiceManager
	Auth          *auth.JWTAuthenticator
	RateLimiter   *ratelimiter.RateLimiter
	Validate      *validator.Validate
	HttpError     *errors.Error
	Monitoring    *monitoring.Monitoring
}

func (api *NotificationApi) Mount() http.Handler {
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
		AllowedOrigins: []string{api.Props.CommonProps.FrontendUrl, "http://localhost"},
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

	// RegisterRoutes registers the notification handler routes
	api.MountNotificationRoutes(router)
	return router
}

// MountNotificationRoutes mounts the notification routes
func (api *NotificationApi) MountNotificationRoutes(r chi.Router) {
	r.Route("/v1/notifications", func(r chi.Router) {
		r.Post("/user-confirmation", api.HandleUserRegisterNotification)
	})
}
