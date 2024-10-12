package auth

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"server/Iternal/Storage"
	"server/Iternal/Storage/models"
	resp "server/Iternal/lib/api/response"
	"time"
)

type Request struct {
	Password    string `json:"password" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Patronymic  string `json:"patronymic"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email" validate:"required,email"`
	Gender      string `json:"gender" validate:"required,oneof=M F "`
}

type Response struct {
	resp.Response
}

type SignUp interface {
	CreateUser(user models.User) error
}

var validate = validator.New()

func SignUpHandler(signUp SignUp, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "SignUpHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to sign up"))
			return
		}

		DateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			log.Error("failed to parse date of birth", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to parse date of birth"))
			return
		}

		err = signUp.CreateUser(models.User{
			Name:           req.Name,
			Patronymic:     req.Patronymic,
			Surname:        req.Surname,
			Email:          req.Email,
			Gender:         req.Gender,
			DateOfBirth:    DateOfBirth,
			PhoneNumber:    req.PhoneNumber,
			HashedPassword: string(hashedPassword),
			Role:           "user",
		})
		if errors.Is(err, Storage.ErrEmailExists) {
			log.Info("email already exists", slog.String("email", req.Email))
			render.JSON(w, r, resp.Error("user with this email already sign-up"))
			return
		}
		if err != nil {
			log.Error("failed to sign up", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to sign up"))
			return
		}

		log.Info("sign up success", slog.String("email", req.Email))

		render.JSON(w, r, resp.OK())
	}
}
