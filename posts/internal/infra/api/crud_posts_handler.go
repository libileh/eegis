package api

import (
	"context"
	"fmt"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/posts/domain"
	"net/http"
)

func (api *PostApi) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("authedUser").(*auth.CtxUser)

	var request domain.PostPayload
	if err := json_utils.ReadJson(w, r, &request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return // Add return
	}
	if err := api.Validate.Struct(request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return // Add return
	}

	newPost, err := domain.NewPost(request, user.ID)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}
	id, customErr := api.Service.PostService.PostRepo.Create(r.Context(), newPost)

	if customErr != nil {
		api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("post creation failed %q", customErr.Error()))
		return // Add return
	}

	if err := json_utils.JsonResponse(w, http.StatusCreated, id); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {

	post, ok := r.Context().Value("post").(*domain.Post)
	if !ok {
		api.HttpError.BadRequestResponse(w, r, "Unable to find post")
	}

	if err := json_utils.JsonResponse(w, http.StatusOK, post); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) getAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, customErr := api.Service.PostService.PostRepo.GetAllPosts(r.Context())
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
		return
	}
	if err := json_utils.JsonResponse(w, http.StatusOK, posts); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := ValidateID(r, "postId")
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	customErr := api.Service.PostService.PostRepo.DeletePost(r.Context(), *id)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
		return
	}
	if err := json_utils.JsonResponse(w, http.StatusNoContent, nil); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, ok := r.Context().Value("post").(*domain.Post)
	if !ok {
		api.HttpError.BadRequestResponse(w, r, "Post not found")
	}
	if err := api.Validate.Struct(post); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}

	var request domain.UpdatePayload
	if err := json_utils.ReadJson(w, r, &request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}
	postToUpdate, err := domain.UpdatePost(&request, post)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}
	result, customErr := api.Service.PostService.PostRepo.UpdatePost(r.Context(), post.ID, post.Version, postToUpdate)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
	}

	if err := json_utils.JsonResponse(w, 200, result); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) getCommentsByPostId(w http.ResponseWriter, r *http.Request) {
	postId, err := ValidateID(r, "postId")
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	comments, customErr := api.Service.PostService.PostRepo.GetCommentsByPostId(r.Context(), *postId)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
		return
	}

	if err := json_utils.JsonResponse(w, http.StatusOK, comments); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) postLoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam, err := ValidateID(r, "postId")
		if err != nil {
			api.HttpError.BadRequestResponse(w, r, err.Error())
			return
		}
		post, customErr := api.Service.PostService.PostRepo.GetPostById(r.Context(), *idParam)
		if customErr != nil {
			api.HttpError.HandleErrorFromDB(w, r, customErr)
			return
		}

		// Store the post in the request context
		ctx := context.WithValue(r.Context(), "post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
