package handlers

//
//import (
//	"github.com/go-chi/render"
//	"github.com/go-playground/validator/v10"
//	"log/slog"
//	"net/http"
//	resp "server/api/lib/response"
//	auth2 "server/internal/features/user/handlers/auth"
//)
//
//type SendConfirmedEmailCodeRequest struct {
//	Email string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
//}
//
//type EmailSendler interface {
//	SendConfirmEmail(code string, email string) error
//}
//
//type CodeGenerator interface {
//	CreateEmailConfirmedCode(email string) (string, error)
//}
//
//// SendConfirmedEmailCodeHandler godoc
//// @Summary Отправка кода подтверждения на email
//// @Tags auth
//// @Description Генерирует код подтверждения и отправляет его на указанный email
//// @Accept json
//// @Produce json
//// @Param request body SendConfirmedEmailCodeRequest true "Email пользователя для подтверждения"
//// @Success 201 {object} UserSignUpResponse "Код подтверждения успешно отправлен"
//// @Failure 400 {object} ValidationErrorResponse "Ошибка валидации или неверный запрос"
//// @Failure 500 {object} InternalServerErrorResponse "Внутренняя ошибка сервера"
//// @Router /auth/email/send-confirm-code [post]
//func SendConfirmedEmailCodeHandler(cg CodeGenerator, es EmailSendler, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
//	log = log.With("op", "SendConfirmedEmailCode")
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		var req SendConfirmedEmailCodeRequest
//
//		if err := render.DecodeJSON(r.Body, &req); err != nil {
//			log.Error("failed to decode request body", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.Error("failed to decode request"))
//			return
//		}
//
//		if err := auth2.validate.Struct(req); err != nil {
//			ValidateErr := err.(validator.ValidationErrors)
//			log.Error("failed to validate request", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.ValidationError(ValidateErr))
//			return
//		}
//
//		code, err := cg.CreateEmailConfirmedCode(req.Email)
//		if err != nil {
//			log.Error("failed to create email confirmed code", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.Error("failed to create email confirmed code"))
//			return
//		}
//
//		if err := es.SendConfirmEmail(code, req.Email); err != nil {
//			log.Error("failed to send email confirmed code", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.Error("failed to send email confirmed code"))
//			return
//		}
//
//		w.WriteHeader(http.StatusCreated)
//		render.JSON(w, r, resp.OK())
//	}
//}
