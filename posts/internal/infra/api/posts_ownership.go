package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/posts/domain"
	"net/http"
)

// checkPostOwnershipHandler handles the GET /v1/posts/{postId}/ownership?userId={userId}&role={role} endpoint.
func (api *PostApi) checkPostOwnershipHandler(w http.ResponseWriter, r *http.Request) {
	// Extract post from context
	post := r.Context().Value("post").(*domain.Post)

	// Extract query parameters
	userQP := r.URL.Query().Get("userId")
	userID, err := uuid.Parse(userQP)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("invalid %s param: %v", userID, err))
		return
	}
	role := r.URL.Query().Get("role")

	// Check if the user owns the post
	if post.UserID == userID {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Call the users service to check role precedence
	allowed, err := api.Service.UserService.CheckRolePrecedence(&userID, role)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
	if !allowed {
		api.HttpError.ForbiddenResponse(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// posts/internal/infra/api/post_api.go

func (api *PostApi) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value("authedUser").(*auth.CtxUser)
		post := r.Context().Value("post").(*domain.Post)

		// Call the ownership endpoint
		resp, err := http.Get(fmt.Sprintf("%s/v1/posts/%s/ownership?userId=%s&role=%s",
			api.Properties.PostsServiceURL,
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
