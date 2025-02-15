package api

import (
	"fmt"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/posts/domain"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	// TODO: Remove hardcoded UUID once user auth is implemented
	tempUserID = "dfb5ac21-0a68-4077-a4e5-72c6417ba82f"
)

func (api *PostApi) createComment(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(tempUserID)
	if err != nil {
		log.Printf("error in userId %s", err)
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	request, err := api.validateCommentPayload(w, r)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("userId error %s", err.Error()))
		return
	}

	comment := &domain.Comment{
		ID:        uuid.New(),
		Content:   request.Content,
		UserID:    userID,
		PostID:    request.PostId,
		CreatedAt: time.Now(),
	}

	id, customErr := api.Service.PostService.CommentRepo.Create(r.Context(), comment)
	if customErr != nil {
		api.HttpError.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("comment creation failed: %v", err))
		return
	}

	err = json_utils.JsonResponse(w, http.StatusCreated, id)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *PostApi) validateCommentPayload(w http.ResponseWriter, r *http.Request) (*domain.CommentPayload, error) {
	var request domain.CommentPayload
	if err := json_utils.ReadJson(w, r, &request); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}

	if err := api.Validate.Struct(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &request, nil
}
