package monitoring

import (
	"database/sql"
	"expvar"
	"runtime"
)

type Monitoring struct {
	Metrics *Metrics
	Health  *Health
}

// Metrics holds the expvar map for monitoring purposes.
type Metrics struct {
	Expvars *expvar.Map
}

// NewMetrics initializes and returns a Metrics instance with version, database stats, and goroutine count.
func NewMetrics(version string, db *sql.DB) *Metrics {
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	return &Metrics{
		Expvars: expvar.NewMap("expvars"),
	}
}
