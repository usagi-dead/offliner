package auth

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"server/internal/lib/api/jwt"
	resp "server/internal/lib/api/response"
	"server/internal/storage/models"
	"time"
)

type SingInRequest struct {
	Email    string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
	Password string `json:"password" validate:"required" example:"SuperPassword123"`
}

type SingInResponse struct {
	Status      string `json:"status" example:"OK"`
	AccessToken string `json:"access_token" example:"asdfasdfahgwea94i5)()(&_.KJFDI{.sadfasdIOSDJ"`
}

type ErrorResponse struct {
	Status string `json:"status" example:"Error"`
	Error  string `json:"error" example:"Invalid Password or Email"`
}

type SignIn interface {
	GetUserByEmail(Email string) (*models.User, error)
}

// SignInHandler godoc
// @Summary Sign In User
// @Tags auth
// @Description Create access and refresh token and return them to the user
// @Accept json
// @Produce json
// @Param user body SingInRequest true "User login details"
// @Success 200 {object} SingInResponse "User successfully signed in"
// @Failure 400 {object} ValidationErrorResponse "Invalid request payload or validation error"
// @Failure 404 {object} ErrorResponse "Invalid Password or Email"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /auth/sign-in [post]
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
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("Invalid Password or Email"))
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(*user.HashedPassword), []byte(req.Password)); err != nil {
			log.Info("Invalid password", err)
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("Invalid Password or Email"))
			return
		}

		if !user.VerifiedEmail {
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, resp.Error("User email is not confirmed"))
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
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.AccessToken(accessToken))
	}
}
