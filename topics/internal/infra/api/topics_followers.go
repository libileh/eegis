package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (api *TopicApi) FollowTopicHandler(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(userIdParam)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, "Invalid user ID")
	}
	topicIdParam := chi.URLParam(r, "topicId")
	topicId, err := uuid.Parse(topicIdParam)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, "Invalid topic ID")
	}

	if customErr := api.Service.TopicRepository.FollowTopic(r.Context(), userId, topicId); customErr != nil {
		api.HttpError.InternalServerError(w, r, customErr.Error())
	}

	w.WriteHeader(http.StatusCreated)
}
