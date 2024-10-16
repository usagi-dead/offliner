package auth

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "server/Iternal/lib/api/response"
)

type EmailConfirmedRequest struct {
	Code  string `json:"code" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type EmailConfirmedCodeGetter interface {
	GetEmailConfirmedCode(email string) (string, error)
}

type EmailConfirmedUpdater interface {
	UpdateEmailStatus(Email string) error
}

func EmailConfirmedHandler(eccg EmailConfirmedCodeGetter, ecu EmailConfirmedUpdater, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	log = log.With("op", "EmailConfirmedHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		var req EmailConfirmedRequest

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

		code, err := eccg.GetEmailConfirmedCode(req.Email)
		if err != nil {
			log.Error("failed to get email confirmed code", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to get email confirmed code"))
			return
		}

		if code != req.Code {
			log.Error("failed to validate email confirmed code", slog.Any("code", code))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to validate email confirmed code"))
			return
		}

		if err := ecu.UpdateEmailStatus(req.Email); err != nil {
			log.Error("failed to update email status", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to update email status"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.OK())
	}
}
