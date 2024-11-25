package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	resp "server/api/lib/response"
	u "server/internal/features/user"
)

// EmailConfirmed godoc
// @Summary Confirmation email address
// @Tags Confirmations
// @Description Validate confirmed code and is it confirmed update email_status
// @Accept json
// @Produce json
// @Param request body EmailConfirmedRequest true "data for confirmed email"
// @Success 200 {object} response.Response "Success email confirmation"
// @Failure 400 {object} response.Response "Error email confirmation"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /confirm/email [put]
func (uc *UserClient) EmailConfirmed(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With("op", "EmailConfirmedHandler")

	var req EmailConfirmedRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("failed to decode request"))
		return
	}

	if err := uc.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(err))
		return
	}

	if err := uc.us.EmailConfirmed(req.Email, req.Code); err != nil {
		switch {
		case errors.Is(err, u.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("email confirmed not found"))
		case errors.Is(err, u.ErrEmailAlreadyConfirmed):
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, resp.Error("email confirmed already exists"))
		default:
			log.Error("failed to confirm email", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.OK())
}
