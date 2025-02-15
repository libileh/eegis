package api

import (
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/request"
	"net/http"
)

func (api *UserApi) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload request.UserRequest
	if err := json_utils.ReadJson(w, r, &payload); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	if err := api.Validate.Struct(&payload); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}
	user := request.MapToUser(payload)

	var passwd domain.Password
	passwd.User = &user
	if err := passwd.Set(payload.Password, api.Service.PasswordService); err != nil {
		api.HttpError.BadRequestResponse(w, r, err.Error())
		return
	}

	var userInvitation domain.UserInvitation
	userInvitation.Token = userInvitation.GenerateToken()

	id, customErr := api.Service.UserRepoService.UserRepo.CreateAndInvite(r.Context(), &user, &userInvitation)
	if customErr != nil {
		api.HttpError.BadRequestResponse(w, r, customErr.Error())
		return
	}

	// Log the registration event
	api.Logger.Infow("User registered", "username", user.Username, "email", user.Email)

	// Send confirmation email using Mailtrap API
	if err := api.Service.NotificationService.SendConfirmationEmail(user.Email, userInvitation.Token); err != nil {
		api.Logger.Errorw("Failed to send confirmation email", "error", err, "email", user.Email)
		// If email fails, rollback user creation
		if err := api.Service.UserRepoService.UserRepo.DeleteUser(r.Context(), *id); err != nil {
			api.Logger.Errorw("Rollback failed: unable to delete user", "error", err, "user_id", *id)
		}
		api.HttpError.InternalServerError(w, r, "Failed to send confirmation email")
		return
	}

	// Activate the user after successful email notification
	if err := api.Service.UserRepoService.UserRepo.ActivateUser(r.Context(), userInvitation.Token); err != nil {
		api.HttpError.HandleErrorFromDB(w, r, err)
		return
	}

	// Return the user ID as the response
	if err := json_utils.JsonResponse(w, http.StatusCreated, *id); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}
}
