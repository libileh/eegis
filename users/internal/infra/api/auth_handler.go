package api

import (
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/json_utils"
	"net/http"
)

func (api *UserApi) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	//parse payload
	var payload auth.AuthRequest
	if err := json_utils.ReadJson(w, r, &payload); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	if err := api.Validate.Struct(&payload); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	//fetch the user from the payload: check if user exist
	user, customErr := api.Service.UserRepoService.UserRepo.GetByEmail(r.Context(), payload.Email)
	if customErr != nil {
		if customErr.ErrType == errors.NotFound {
			api.HttpError.UnauthorizedErrorResponse(w, r, customErr.Error())
		} else {
			api.HttpError.InternalServerError(w, r, customErr.Error())
		}
		return
	}

	claims, err := api.Service.MapToJWTClaims(user.ID, user.RoleID, api.Properties.CommonProps.AuthProps)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
	token, err := api.Auth.GenerateToken(claims)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
	if err := json_utils.JsonResponse(w, http.StatusCreated, token); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

// RefreshTokenHandler handles refresh token requests via GET method.
// It expects a query parameter "user_id", and generates a new token for the user.
// In production, this endpoint should be secured (e.g., by validating an existing refresh token).
func (api *UserApi) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the query parameters.
	userIdParam := r.URL.Query().Get("user_id")
	if userIdParam == "" {
		api.HttpError.BadRequestResponse(w, r, "user_id query parameter is required")
		return
	}

	userId, err := uuid.Parse(userIdParam)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, "invalid user_id parameter")
		return
	}

	// Retrieve the user from the repository.
	// Note: This assumes the existence of a GetByID method on the user repository.
	user, customErr := api.Service.UserRepoService.UserRepo.GetById(r.Context(), userId)
	if customErr != nil {
		if customErr.ErrType == errors.NotFound {
			api.HttpError.UnauthorizedErrorResponse(w, r, customErr.Error())
		} else {
			api.HttpError.InternalServerError(w, r, customErr.Error())
		}
		return
	}

	claims, err := api.Service.MapToJWTClaims(user.ID, user.RoleID, api.Properties.CommonProps.AuthProps)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}

	// Generate a new token.
	newToken, err := api.Auth.GenerateToken(claims)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}

	// Return the token with a standard OK response.
	if err := json_utils.JsonResponse(w, http.StatusOK, newToken); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
}
