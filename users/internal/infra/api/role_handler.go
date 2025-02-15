package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/users/domain"
	"net/http"
)

// checkRolePrecedenceHandler handles the GET /v1/users/role-precedence endpoint.
func (api *UserApi) checkRolePrecedenceHandler(w http.ResponseWriter, r *http.Request) {
	userQP := r.URL.Query().Get("userId")
	roleName := r.URL.Query().Get("role")

	userID, err := uuid.Parse(userQP)
	if err != nil {
		api.HttpError.BadRequestResponse(w, r, fmt.Sprintf("invalid %s param: %v", userID, err))
	}
	// Fetch the user and their role
	user, err := api.Service.UserRepoService.UserRepo.GetById(r.Context(), userID)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}

	// Check role precedence
	allowed, err := api.checkRolePrecedence(r.Context(), user, roleName)
	if err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
		return
	}

	// Return the result
	if err := json_utils.JsonResponse(w, http.StatusOK, map[string]bool{"allowed": allowed}); err != nil {
		api.HttpError.InternalServerError(w, r, err.Error())
	}

}

func (api *UserApi) checkRolePrecedence(ctx context.Context, user *domain.User, roleName string) (bool, *errors.CustomError) {
	role, err := api.Service.UserRepoService.RoleRepo.GetByRoleName(ctx, roleName)
	if err != nil {
		return false, err
	}

	userRole, err := api.Service.UserRepoService.RoleRepo.GetById(ctx, user.RoleID)
	if err != nil {
		return false, err
	}
	return userRole.Level >= role.Level, nil
}
