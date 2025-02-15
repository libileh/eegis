package api

import (
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/notifications/domain"
	"net/http"
)

// HandleNotification handles POST requests to send notifications
func (api *NotificationApi) HandleNotification(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request domain.Notification
	if err := json_utils.ReadJson(w, r, &request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return // Add return
	}

	// Send the notification
	if err := api.ServerManager.NotificationService.SendEmail(
		request.Type,
		request.Recipient,
		request.Content,
	); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}
