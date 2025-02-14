package ratelimiter

import (
	"sync"
	"time"
)

// FixedWindowAlgo implements a fixed window rate limiting algorithm
type FixedWindowAlgo struct {
	sync.RWMutex                // Mutex for thread-safe operations
	clients      map[string]int // Maps IP addresses to request counts
	limit        int            // Maximum requests allowed per window
	window       time.Duration  // Time window duration
}

// NewFixedWindowLimiter creates a new rate limiter with specified limit and window size
func NewFixedWindowLimiter(limit int, windowSize time.Duration) *FixedWindowAlgo {
	return &FixedWindowAlgo{
		clients: make(map[string]int),
		limit:   limit,
		window:  windowSize,
	}
}

// Allow checks if a request from the given IP is allowed
// Returns:
//   - bool: true if request is allowed, false if rate limit exceeded
//   - time.Duration: wait time before next allowed request (0 if allowed)
func (rl *FixedWindowAlgo) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exits := rl.clients[ip]
	rl.RUnlock()
	if !exits || count < rl.limit {
		rl.Lock()
		if !exits {
			go rl.resetCount(ip) // Start window timer for new IPs
		}
		rl.clients[ip]++
		rl.Unlock()
		return true, 0
	}
	return false, rl.window
}

// resetCount removes the IP from tracking after window duration
// Used internally to reset counters after window expiration
func (rl *FixedWindowAlgo) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}
