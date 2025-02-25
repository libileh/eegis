package api

import (
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/users/domain"
	"net/http"
)

func (api *UserApi) followTopicHandler(w http.ResponseWriter, r *http.Request) {
	ctxUser := r.Context().Value("authedUser").(*auth.CtxUser)
	topicId, err := ValidateID(r, "topicId")

	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	if customErr := api.Service.FollowTopic(ctxUser, *topicId); customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
		return
	}
	if err := json_utils.JsonResponse(w, http.StatusOK, nil); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *UserApi) unfollowTopicHandler(w http.ResponseWriter, r *http.Request) {
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
