package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	ck "server/api/lib/contextKeys"
	resp "server/api/lib/response"
	u "server/internal/features/user"
)

// DeleteUserHandler deletes a user from the system.
//
// @Summary      Delete user
// @Description  Deletes the currently authenticated user based on the user ID retrieved from the request context.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response "User successfully deleted"
// @Failure      401  {object}  response.Response "Unauthorized or invalid user ID"
// @Failure      500  {object}  response.Response "Internal server error"
// @Router       /user/delete [delete]
func (uc *UserClient) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With("op", "DeleteUserHandler")

	userId, ok := r.Context().Value(ck.UserIDKey).(int64)
	if !ok {
		log.Error("can't get userId from context")
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("invalid user id"))
		return
	}

	if err := uc.us.DeleteUser(userId); err != nil {
		log.Error(err.Error())
		switch {
		case errors.Is(err, u.ErrInternal):
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, u.ErrInternal.Error())
		default:
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(u.ErrInternal.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StatusOK)
	return
}
