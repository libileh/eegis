package monitoring

import (
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/json_utils"
	"net/http"
)

type Health struct {
	Version   string
	Env       string
	HttpError *errors.Error
}

func (health *Health) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     health.Env,
		"version": health.Version,
	}

	if err := json_utils.WriteJSON(w, http.StatusOK, data); err != nil {
		health.HttpError.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}

func NewHealth(version, env string) *Health {
	return &Health{Version: version, Env: env}
}
