package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "server/api/lib/response"
	u "server/internal/features/user"
)

// RefreshToken
// @Summary Refresh Access Token
// @Tags Authentication
// @Description Refreshes the access token using the provided refresh token from cookies.
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Successfully refreshed access token"
// @Failure 401 {object} response.Response "Invalid, missing or expired refresh token"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/refresh-token [post]
func (uc *UserClient) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With(slog.String("op", "RefreshTokenHandler"))

	AccessToken, err := uc.us.RefreshToken(r)
	if err != nil {
		switch {
		case errors.Is(err, u.ErrNoRefreshToken):
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error(err.Error()))
		case errors.Is(err, u.ErrInvalidToken):
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error(err.Error()))
		case errors.Is(err, u.ErrExpiredToken):
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error(err.Error()))
		case errors.Is(err, u.ErrUserNotFound):
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, resp.Error(err.Error()))
		default:
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(u.ErrInternal.Error()))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.AccessToken(AccessToken))
}
