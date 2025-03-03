package api

import (
	"net/http"

	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/notifications/domain"
)

// HandleUserRegisterNotification handles POST requests to send user verification emails.
func (api *NotificationApi) HandleUserRegisterNotification(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request domain.Notification
	if err := json_utils.ReadJson(w, r, &request); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	// Send the verification email
	if err := api.ServerManager.NotificationService.SendUserVerificationEmail(
		request.Recipient,
		request.Content,
	); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}

// HandleSendPostReviewNotification handles the notification trigger for post review events.
// HandleSendPostReviewNotification handles the notification trigger for post review events.
func (api *NotificationApi) HandleSendPostReviewNotification(w http.ResponseWriter, r *http.Request) {
	// Start the notification trigger for post review events using the dedicated service function.
	if err := api.ServerManager.NotificationService.StartReviewPostNotificationTrigger(); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}
	// Return success response
	w.WriteHeader(http.StatusOK)
}
