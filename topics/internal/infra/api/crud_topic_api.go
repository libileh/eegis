package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/topics/domain"
	"github.com/libileh/eegis/topics/internal/infra/request"
	"net/http"
)

func (api *TopicApi) CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var topicPayload request.TopicPayload
	if err := json_utils.ReadJson(w, r, &topicPayload); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}
	if err := api.Validator.Struct(topicPayload); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}
	newTopic, err := domain.NewTopic(topicPayload)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}
	id, customErr := api.Service.TopicRepository.Create(r.Context(), newTopic)
	if customErr != nil {
		api.HttpError.InternalServerError(w, r, customErr.Error())
	}
	if err := json_utils.JsonResponse(w, http.StatusCreated, id); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *TopicApi) GetTopicByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	topic, customErr := api.Service.TopicRepository.GetTopicByName(r.Context(), name)
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
	}

	if err := json_utils.JsonResponse(w, http.StatusOK, topic); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *TopicApi) GetUserFollowedTopicsHandler(w http.ResponseWriter, r *http.Request) {
	paramID := chi.URLParam(r, "followerId")
	followerId, err := uuid.Parse(paramID)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
	}

	followerTopics, customErr := api.Service.GetFollowerTopic(r.Context(), followerId)
	api.HandlerError(w, r, customErr)

	if err := json_utils.JsonResponse(w, http.StatusOK, followerTopics); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}

func (api *TopicApi) GetAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	topics, customErr := api.Service.TopicRepository.GetAllTopics(r.Context())
	if customErr != nil {
		api.HttpError.HandleErrorFromDB(w, r, customErr)
	}
	if err := json_utils.JsonResponse(w, http.StatusOK, topics); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}
