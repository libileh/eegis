package ratelimiter

import (
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/properties"
	"net/http"
	"time"
)

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type RateLimiter struct {
	Properties      *properties.RateLimiterProperties
	RateLimiterAlgo *FixedWindowAlgo
	Error           *errors.Error
}

func (rl *RateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.Properties.Enabled {
			if allow, retryAfter := rl.RateLimiterAlgo.Allow(r.RemoteAddr); !allow {
				rl.Error.RateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
