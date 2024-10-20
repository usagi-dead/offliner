package auth

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	resp "server/internal/lib/api/response"
	"server/internal/storage"
)

type UserSignUpRequest struct {
	Email    string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
	Password string `json:"password" validate:"required" example:"SuperPassword123"`
}

type UserSignUpResponse struct {
	Status string `json:"status" example:"OK"`
}

type ValidationErrorResponse struct {
	Status string `json:"status" example:"Error"`
	Error  string `json:"error" example:"field Email is not a valid Email, field Password is a required field"`
}

type ConflictErrorResponse struct {
	Status string `json:"status" example:"Error"`
	Error  string `json:"error" example:"user with this email already sign-up"`
}

type InternalServerErrorResponse struct {
	Status string `json:"status" example:"Error"`
	Error  string `json:"error" example:"failed to sign up"`
}

type SignUp interface {
	CreateUser(email string, hashedPassword string) error
}

var validate = validator.New()

// SignUpHandler godoc
// @Summary Sign Up User
// @Tags auth
// @Description Creates a new user account with the provided details.
// @Accept json
// @Produce json
// @Param user body UserSignUpRequest true "User registration details"
// @Success 201 {object} UserSignUpResponse "User successfully registered"
// @Failure 400 {object} ValidationErrorResponse "Invalid request payload or validation error"
// @Failure 409 {object} ConflictErrorResponse "User with this email already exists"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /auth/sign-up [post]
func SignUpHandler(signUp SignUp, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "SignUpHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UserSignUpRequest

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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to sign up"))
			return
		}

		err = signUp.CreateUser(req.Email, string(hashedPassword))

		if err != nil {
			if errors.Is(err, storage.ErrEmailExists) {
				log.Info("email already exists", slog.String("email", req.Email))
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, resp.Error("user with this email already sign-up"))
				return
			}
			log.Error("failed to sign up", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to sign up"))
			return
		}

		log.Info("sign up success", slog.String("email", req.Email))
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, resp.OK())
	}
}
