package errors

import (
	"github.com/libileh/eegis/common/json_utils"
	"go.uber.org/zap"
	"net/http"
)

type Error struct {
	Logger *zap.SugaredLogger
}

func (err *Error) WriteJSONError(w http.ResponseWriter, status int, message string) error {
	err.Logger.Errorf("Error response: status=%d, message=%s", status, message)
	return json_utils.WriteJSON(w, status, message)
}

func (err *Error) HandleErrorFromDB(w http.ResponseWriter, r *http.Request, customErr *CustomError) {
	switch customErr.ErrType {
	case NotFound:
		err.NotFoundResponse(w, r, customErr.Message)
	default:
		err.InternalServerError(w, r, customErr.Message)
	}
}

func (err *Error) InternalServerError(w http.ResponseWriter, r *http.Request, message string) {
	err.WriteJSONError(w, http.StatusInternalServerError, message)
}

func (err *Error) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	err.Logger.Warnf("forbiden Request", "path", r.URL.Path, "method", r.Method)
	err.WriteJSONError(w, http.StatusForbidden, "user is not authorized to perform this operation")
}

func (err *Error) BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	err.Logger.Errorf("Bad Request", "path", r.URL.Path, "error", message)
	err.WriteJSONError(w, http.StatusBadRequest, message)
}

func (err *Error) ConflictResponse(w http.ResponseWriter, r *http.Request, message string) {
	err.WriteJSONError(w, http.StatusConflict, message)
}

func (err *Error) NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	err.WriteJSONError(w, http.StatusNotFound, message)
}
func (err *Error) UnauthorizedTokenResponse(w http.ResponseWriter, r *http.Request, message string) {
	err.Logger.Warnf("unauthorized token error", "method", r.Method, "path", r.URL.Path, "error", message)
	err.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (err *Error) UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	err.Logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", message)
	w.Header().Set("WWW-Authenticate", `"Basic realm=Restricted" charset="UTF-8"`)
	err.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (err *Error) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	w.Header().Set("Retry-After", retryAfter)
	err.WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
