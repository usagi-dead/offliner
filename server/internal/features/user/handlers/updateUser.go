package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"mime/multipart"
	"net/http"
	ck "server/api/lib/contextKeys"
	resp "server/api/lib/response"
	u "server/internal/features/user"
)

// UpdateUserHandler
// @Security ApiKeyAuth
// @Summary Update User Profile
// @Description Updates user profile details including name, surname, patronymic, date of birth, phone number, gender, and avatar.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param json formData string true "User profile data (JSON string)"
// @Param avatar formData file false "Avatar image file (optional)"
// @Success 200 {object} response.Response "User profile updated successfully"
// @Failure 400 {object} response.Response "Invalid request or validation error"
// @Failure 401 {object} response.Response "Unauthorized user"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/edit [put]
func (uc *UserClient) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log := uc.log.With("op", "UpdateUserHandler")

	userId, ok := r.Context().Value(ck.UserIDKey).(int64)
	if !ok {
		log.Error("can't get userId from context")
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, resp.Error("invalid user id"))
		return
	}

	var req UpdateUserRequest

	if err := r.ParseMultipartForm(1 << 20); err != nil { // Limit to 1MB
		log.Error("failed to parse multipart form", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error(u.ErrSizeAvatar.Error()+" or invalid request"))
		return
	}

	jsonData := r.Form.Get("json")
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		log.Error("failed to decode JSON part", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("invalid JSON part"))
		return
	}

	if err := uc.validate.Struct(req); err != nil {
		log.Error("failed to validate request", err)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(err))
		return
	}

	var avatarFile *multipart.File
	file, _, _ := r.FormFile("avatar")
	if file != nil {
		defer file.Close()
	}

	user := &u.User{
		UserId:      userId,
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		DateOfBirth: req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
	}

	if err := uc.us.UpdateUser(avatarFile, user, req.ResetAvatar); err != nil {
		log.Error("failed to update user", err)
		switch {
		case errors.Is(err, u.ErrUserNotFound):
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(u.ErrUserNotFound.Error()))
		case errors.Is(err, u.ErrInvalidResolutionAvatar) || errors.Is(err, u.ErrInvalidTypeAvatar):
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(u.ErrInternal.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.OK())
	return
}
