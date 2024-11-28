package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "server/api/lib/response"
	"server/internal/features/user"
)

// SignUp
// @Summary User SignUp
// @Tags Authentication
// @Description Registers a new user with the provided email and password.
// @Accept json
// @Produce json
// @Param user body handlers.UserSignUpRequest true "User registration details"
// @Success 201 {object} response.Response "User successfully created"
// @Failure 400 {object} response.Response "Validation error or invalid request payload"
// @Failure 409 {object} response.Response "Conflict - User with this email already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/sign-up [post]
func (uc *UserClient) SignUp(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With(slog.String("op", "SignUpHandler"))

	var req UserSignUpRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("failed to decode request"))
		return
	}

	if err := uc.validate.Struct(req); err != nil {
		log.Error("failed to validate request", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(err))
		return
	}

	if err := uc.us.SignUp(req.Email, req.Password); err != nil {
		switch {
		case errors.Is(err, user.ErrEmailExists):
			log.Info("email already exists", slog.String("email", req.Email))
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, resp.Error("user with this email already exists"))
		default:
			log.Error("failed to sign up user", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.OK())
}
