package auth

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"server/Iternal/lib/api/jwt"
	resp "server/Iternal/lib/api/response"
	"server/Iternal/storage/models"
	"time"
)

type SingInRequest struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type SingInResponse struct {
	AccessToken string `json:"access_token"`
	resp.Response
}

type SignIn interface {
	GetUserByEmail(Email string) (*models.User, error)
}

func SignInHandler(signIn SignIn, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "SignInHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SingInRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("body", req))

		if err := validate.Struct(req); err != nil {
			ValidateErr := err.(validator.ValidationErrors)
			log.Error("failed to validate request", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(ValidateErr))
			return
		}

		user, err := signIn.GetUserByEmail(req.Email)
		if err != nil {
			log.Error("User not found", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid Password or Email"))
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
			log.Error("Invalid password", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Invalid Password or Email"))
			return
		}

		accessToken, err := jwt.GenerateAccessToken(user.UserId, user.Role)
		if err != nil {
			log.Error("failed to generate access token", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to generate access token"))
			return
		}

		refreshToken, err := jwt.GenerateRefreshToken(user.UserId)
		if err != nil {
			log.Error("failed to generate refresh token", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to generate refresh token"))
			return
		}

		log.Info("sign in success", slog.String("email", req.Email))

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(15 * 24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")

		render.JSON(w, r, resp.AccessToken(accessToken))
	}
}
