package profile

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

const avatarDir = "./avatars"

// AvatarHandler serves the avatar image for a user.
// @Security ApiKeyAuth
// @Summary Get User Avatar
// @Description Retrieves the avatar image based on the filename provided in the query parameter.
// @Tags profile
// @Accept  json
// @Produce  json
// @Param filename query string true "Filename of the avatar image"
// @Success 200 {file} file "Avatar image"
// @Router /user/avatar [get]
func AvatarHandler(log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	log = log.With("op", "avatarHandler")

	return func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("filename")
		if fileName == "" {
			http.Error(w, "filename required", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(avatarDir, fileName)
		log.Info("Serving file", "path", filePath)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, filePath)
	}
}
