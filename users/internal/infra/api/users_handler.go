package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/users/domain"
	"log"
	"net/http"
)

func (api *UserApi) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := ValidateID(r, "userId")
	if err != nil {
		log.Printf("bad userId %v", err)
		api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("bad userId %v", err))
		return
	}

	user, customErr := api.Service.UserRepoService.UserRepo.GetById(r.Context(), *id)
	if customErr != nil {
		log.Printf("userId not found: %v", err)
		api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("userId not found: %v", err))
		return
	}

	err = json_utils.JsonResponse(w, http.StatusOK, user)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
}

func (api *UserApi) followUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := r.Context().Value("authedUser").(*auth.CtxUser)

	followedId, err := ValidateID(r, "userId")
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	follow, customErr := api.Service.UserRepoService.FollowerRpo.FollowUser(r.Context(),
		&domain.Follower{FollowerId: follower.ID, UserId: *followedId})
	if customErr != nil {
		api.HttpError.BadRequestResponse(w, r, customErr.Error())
		return
	}
	if err := json_utils.JsonResponse(w, http.StatusOK, follow); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *UserApi) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := r.Context().Value("authedUser").(*auth.CtxUser)

	followedId, err := ValidateID(r, "userId")
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	customErr := api.Service.UserRepoService.FollowerRpo.UnfollowUser(
		r.Context(), &domain.Follower{FollowerId: follower.ID, UserId: *followedId})
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
		return
	}
	if err := json_utils.JsonResponse(w, http.StatusNoContent, nil); err != nil {
		return
	}

}

func (api *UserApi) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	customErr := api.Service.UserRepoService.UserRepo.ActivateUser(r.Context(), token)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
	}

	if err := json_utils.JsonResponse(w, http.StatusOK, nil); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

//userLoaderMiddleware was user before authentication middleware

func (api *UserApi) userLoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := ValidateID(r, "userId")

		if err != nil {
			log.Printf("bad userId %v", err)
			api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("bad userId %v", err))
			return
		}
		user, customErr := api.Service.UserRepoService.UserRepo.GetById(r.Context(), *id)
		if customErr != nil {
			log.Printf("userId not found: %v", err)
			api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("userId not found: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *UserApi) getUser(ctx context.Context, userId uuid.UUID) (*domain.User, *errors.CustomError) {
	if !api.Properties.Redis.Enabled {
		return api.Service.UserRepoService.UserRepo.GetById(ctx, userId)
	}
	user, customErr := api.Service.UserCacheService.UserCache.Get(ctx, userId)
	if customErr != nil {
		return nil, customErr
	}
	if user != nil {
		return user, nil
	}

	user, customErr = api.Service.UserRepoService.UserRepo.GetById(ctx, userId)
	if customErr != nil {
		return nil, customErr
	}
	customErr = api.Service.UserCacheService.UserCache.Set(ctx, user)
	if customErr != nil {
		return nil, customErr
	}
	return user, nil
}
