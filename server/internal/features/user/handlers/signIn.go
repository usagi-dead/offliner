package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "server/api/lib/response"
	"server/internal/features/user"
	"time"
)

// SignIn
// @Summary User SignIn
// @Tags Authentication
// @Description Create access and refresh token and return them to the user
// @Accept json
// @Produce json
// @Param user body handlers.UserSingInRequest true "User login details"
// @Success 200 {object} response.Response "User successfully signed in"
// @Failure 400 {object} response.Response "Invalid request payload or validation error"
// @Failure 401 {object} response.Response "Invalid Password or Email"
// @Failure 403 {object} response.Response "User email is not confirmed"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/sign-in [post]
func (uc *UserClient) SignIn(w http.ResponseWriter, r *http.Request) {

	log := uc.log.With(slog.String("op", "SignInHandler"))

	var req UserSingInRequest

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

	AccessToken, RefreshToken, err := uc.us.SignIn(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("failed email or password"))
		case errors.Is(err, user.ErrEmailNotConfirmed):
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("email not confirmed"))
		default:
			log.Error("failed to sign in", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    RefreshToken,
		Expires:  time.Now().Add(15 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.AccessToken(AccessToken))
}
