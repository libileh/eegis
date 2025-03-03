package api

import (
	"fmt"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/posts/domain"
	"net/http"
)

// posts/internal/infra/api/post_api.go

func (api *PostApi) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value("authedUser").(*auth.CtxUser)
		post := r.Context().Value("post").(*domain.Post)

		// Check if the user owns the post
		if post.UserID == ctxUser.ID {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the users service to check role precedence
		allowed, err := api.Service.UserService.CheckRolePrecedence(ctxUser)
		if err != nil {
			api.HttpError.InternalServerError(w, r, err.Error())
			return
		}
		if !allowed {
			api.HttpError.ForbiddenResponse(w, r)
			return
		}

		// Call the ownership endpoint
		resp, err := http.Get(fmt.Sprintf("%s/v1/users/%s/ownership?userId=%s&role=%s",
			api.Properties.UsersServiceURL,
			post.ID,
			ctxUser.ID,
			role))
		if err != nil {
			api.HttpError.InternalServerError(w, r, err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			api.HttpError.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}
