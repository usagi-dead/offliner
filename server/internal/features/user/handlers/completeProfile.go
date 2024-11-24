package handlers

//
//import (
//	"encoding/json"
//	"errors"
//	"github.com/go-chi/render"
//	"github.com/go-playground/validator/v10"
//	"log/slog"
//	"net/http"
//	"server/api/lib/CustomErrors"
//	"server/api/lib/avatarMenager"
//	ck "server/api/lib/contextKeys"
//	resp "server/api/lib/response"
//	"server/internal/features/user/data"
//	"time"
//)
//
//type CompleteProfileRequest struct {
//	Surname     *string    `json:"surname,omitempty" validate:"omitempty,min=1,max=50" example:"John"`
//	Name        *string    `json:"name,omitempty" validate:"omitempty,min=1,max=50" example:"Doe"`
//	Patronymic  *string    `json:"patronymic,omitempty" validate:"omitempty,min=1,max=50" example:"Smith"`
//	DateOfBirth *time.Time `json:"date_of_birth,omitempty" validate:"omitempty,lte" example:"1985-04-12T00:00:00Z"`
//	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+1234567890"`
//	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female" example:"male"`
//}
//
//type CompleteProfileResponse struct {
//	Status  string
//	Message string
//}
//
//type DataSource interface {
//	GetUserById(user *data.User) error
//}
//
//var validate = validator.New()
//
//// CompleteProfileHandler updates the user's profile information including optional avatar upload.
//// @Security ApiKeyAuth
//// @Summary Update User Profile
//// @Description Updates user profile details such as name, surname, patronymic, date of birth, phone number, gender, and avatar.
//// @Tags profile
//// @Accept multipart/form-data
//// @Produce json
//// @Param json formData string true "User profile data JSON"
//// @Param avatar formData file false "Avatar image file"
//// @Success 201 {object} CompleteProfileResponse "User updated successfully"
//// @Router /user/complete-profile [put]
//func CompleteProfileHandler(ds DataSource, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
//	log = log.With("op", "CompleteProfileHandler")
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		userId, ok := r.Context().Value(ck.UserIDKey).(int64)
//		if !ok {
//			http.Error(w, "user_id not found or invalid", http.StatusBadRequest)
//			return
//		}
//
//		var req CompleteProfileRequest
//
//		if err := r.ParseMultipartForm(1 << 20); err != nil { // Limit to 1MB
//			log.Error("failed to parse multipart form", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.Error("failed to parse form"))
//			return
//		}
//
//		jsonData := r.Form.Get("json")
//		if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
//			log.Error("failed to decode JSON part", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.Error("failed to decode JSON"))
//			return
//		}
//
//		if err := validate.Struct(req); err != nil {
//			ValidateErr := err.(validator.ValidationErrors)
//			log.Error("failed to validate request", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.ValidationError(ValidateErr))
//			return
//		}
//
//		file, _, err := r.FormFile("avatar")
//		if err != nil {
//			log.Error("failed to get avatar file", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, resp.Error("failed to get avatar"))
//			return
//		}
//
//		defer file.Close()
//
//		avatarUrl, err := avatarMenager.SaveAvatar(&file)
//		if err != nil {
//			log.Error("failed to save avatar", err)
//
//			// Determine if the error is a user error or a server error
//			var userError bool
//			var statusCode int
//			var errorMessage string
//
//			switch {
//			case errors.Is(err, CustomErrors.ErrUnsupportedImageFormat),
//				errors.Is(err, CustomErrors.ErrAnimatedGIFNotSupported),
//				errors.Is(err, CustomErrors.ErrImageMustBeSquare):
//				userError = true
//				statusCode = http.StatusBadRequest
//				errorMessage = err.Error() // User-facing error message
//
//			case errors.Is(err, CustomErrors.ErrFileReadError),
//				errors.Is(err, CustomErrors.ErrImageDecodingError),
//				errors.Is(err, CustomErrors.ErrFileSaveError),
//				errors.Is(err, CustomErrors.ErrWebPEncodingError):
//				userError = false
//				statusCode = http.StatusInternalServerError
//				errorMessage = "internal server error" // Generic server error message
//
//			default:
//				userError = false
//				statusCode = http.StatusInternalServerError
//				errorMessage = "unknown error occurred" // Fallback for unexpected errors
//			}
//
//			// Respond to the client
//			w.WriteHeader(statusCode)
//			if userError {
//				render.JSON(w, r, resp.Error(errorMessage))
//			} else {
//				render.JSON(w, r, resp.Error("internal server error"))
//			}
//			return
//		}
//
//		user := &data.User{
//			UserId:      userId,
//			Name:        req.Name,
//			Surname:     req.Surname,
//			Patronymic:  req.Patronymic,
//			DateOfBirth: req.DateOfBirth,
//			PhoneNumber: req.PhoneNumber,
//			AvatarUrl:   &avatarUrl,
//			Gender:      req.Gender,
//		}
//
//		if err := ds.GetUserById(user); err != nil {
//			if errors.As(err, &CustomErrors.ErrEmailNotExists) {
//				log.Error("failed to update user", err)
//				w.WriteHeader(http.StatusBadRequest)
//				render.JSON(w, r, resp.Error("user not exists"))
//				return
//			}
//			log.Error("failed to update user", err)
//			w.WriteHeader(http.StatusInternalServerError)
//			render.JSON(w, r, resp.Error("failed to update user"))
//			return
//		}
//
//		w.WriteHeader(http.StatusCreated)
//		render.JSON(w, r, CompleteProfileResponse{
//			Status:  "OK",
//			Message: "user updated successfully",
//		})
//		return
//	}
//}
