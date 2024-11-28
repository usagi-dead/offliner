package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	resp "server/api/lib/response"
	u "server/internal/features/user"
)

// SendConfirmedEmailCode
// @Summary Send code for confirmation email
// @Tags Confirmations
// @Description Generate code for confirmation email and send this to email. This endpoint have rate 1 req in 1 min
// @Accept json
// @Produce json
// @Param request body SendConfirmedEmailCodeRequest true "Email пользователя для подтверждения"
// @Success 201 {object} response.Response "Код подтверждения успешно отправлен"
// @Failure 400 {object} response.Response "Ошибка валидации или неверный запрос"
// @Failure 500 {object} response.Response "Внутренняя ошибка сервера"
// @Router /confirm/send-email-code [post]
func (uc *UserClient) SendConfirmedEmailCode(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With("op", "SendConfirmedEmailCode")

	var req SendConfirmedEmailCodeRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("failed to decode request"))
		return
	}

	if err := uc.validate.Struct(req); err != nil {
		log.Info("failed to validate request data", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(err))
		return
	}

	if err := uc.us.SendEmailForConfirmed(req.Email); err != nil {

		switch {
		case errors.Is(err, u.ErrEmailAlreadyConfirmed):
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(u.ErrEmailAlreadyConfirmed.Error()))
		case errors.Is(err, u.ErrUserNotFound):
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(u.ErrUserNotFound.Error()))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.OK())
	return
}
