package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	ck "server/api/lib/contextKeys"
	resp "server/api/lib/response"
	u "server/internal/features/user"
)

// GetUserHandler
// @Security ApiKeyAuth
// @Summary Get User Profile
// @Description Retrieves the profile details of the authenticated user.
// @Tags User
// @Produce json
// @Success 200 {object} response.Response "User profile retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized user"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/me [get]
func (uc *UserClient) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With("op", "GetUserHandler")

	userId, ok := r.Context().Value(ck.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("invalid user id"))
		return
	}

	user, err := uc.us.GetUser(userId)
	if err != nil {
		log.Error(err.Error())
		switch {
		case errors.Is(err, u.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("user not found"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(u.ErrInternal.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.UserProfile(user))
}
