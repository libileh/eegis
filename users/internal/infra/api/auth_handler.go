package api

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/json_utils"

	"net/http"
	"time"
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

	//generate token: add claims
	ctxUser, err := auth.MapToCtxUser(float64(user.RoleID))
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": ctxUser,
		"exp":  time.Now().Add(api.Properties.CommonProps.AuthProps.Token.Exp).Unix(),
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"iss":  api.Properties.CommonProps.AuthProps.Issuer,
		"aud":  api.Properties.CommonProps.AuthProps.Audience,
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
