package api

import (
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/posts/domain"
	"net/http"
)

func (api *PostApi) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.CtxUser)
	if !ok {
		api.HttpError.BadRequestResponse(w, r, "user not found")
	}

	fp := domain.Paginated{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   make([]string, 0),
		Search: "",
	}
	paginated := fp.GetFeedPaginationFilter(r)

	if err := api.Validate.Struct(paginated); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	userFeed, customErr := api.Service.PostService.FeedRepo.GetFeed(r.Context(), user.ID, paginated)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
	}
	err := json_utils.JsonResponse(w, http.StatusOK, userFeed)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
}
