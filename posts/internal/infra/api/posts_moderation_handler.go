package api

import (
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/posts/internal/infra/requests"
	"net/http"
)

func (api *PostApi) PostModerationHandler(w http.ResponseWriter, r *http.Request) {
	// Your implementation here
	ctxUser := r.Context().Value("authedUser").(*auth.CtxUser)
	postId, err := ValidateID(r, "postId")
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	var request requests.ReviewPostPayload
	if err := json_utils.ReadJson(w, r, &request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}

	if err := api.Validate.Struct(request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}

	postReview, customErr := api.Service.PostService.ReviewPost(r.Context(), postId, ctxUser, request.Status)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
		return
	}

	if err := json_utils.JsonResponse(w, http.StatusOK, postReview); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}
