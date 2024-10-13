package auth

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"server/Iternal/Storage/models"
	"server/Iternal/lib/api/jwt"
	resp "server/Iternal/lib/api/response"
)

type refreshAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	resp.Response
}

type GetUser interface {
	GetUserById(UserId int64) (*models.User, error)
}

func RefreshTokenHandler(getUser GetUser, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "RefreshTokenHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			log.Error("failed extract access token: %v", r)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ValidateJWT(refreshToken.Value)
		if err != nil {
			log.Error("invalid token: %v", err.Error())
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		user, err := getUser.GetUserById(claims.UserId)
		if err != nil {
			log.Error("failed find user by id: %v", err.Error())
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		accessToken, err := jwt.GenerateAccessToken(user.UserId, user.Role)
		if err != nil {
			log.Error("failed to generate access token", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to generate access token"))
			return
		}

		w.Header().Set("Content-Type", "application/json")

		render.JSON(w, r, resp.AccessToken(accessToken))
	}
}
