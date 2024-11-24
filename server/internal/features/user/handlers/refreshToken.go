package handlers

//
//import (
//	"github.com/go-chi/chi/v5/middleware"
//	"github.com/go-chi/render"
//	"log/slog"
//	"net/http"
//	"server/api/lib/jwt"
//	resp "server/api/lib/response"
//	"server/internal/features/user/data"
//)
//
//type refreshAccessTokenResponse struct {
//	AccessToken string `json:"access_token"`
//}
//
//type GetUser interface {
//	GetUserById(UserId int64) (*data.User, error)
//}
//
//type UnauthorizedResponse struct {
//	Status string `json:"status" example:"error"`
//	Error  string `json:"error" example:"unauthorized"`
//}
//
//// RefreshTokenHandler godoc
//// @Summary Refresh Access Token
//// @Tags auth
//// @Description Refreshes the access token using the provided refresh token from cookies.
//// @Accept json
//// @Produce json
//// @Success 200 {object} SingInResponse "Successfully refreshed access token"
//// @Failure 401 {object} UnauthorizedResponse "Invalid or missing refresh token"
//// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
//// @Router /auth/refresh-token [post]
//func RefreshTokenHandler(getUser GetUser, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op string = "RefreshTokenHandler"
//
//		log = log.With(
//			slog.String("op", op),
//			slog.String("request_id", middleware.GetReqID(r.Context())),
//		)
//
//		refreshToken, err := r.Cookie("refresh_token")
//		if err != nil {
//			log.Error("failed extract access token: %v", r)
//			http.Error(w, "token missed", http.StatusUnauthorized)
//			return
//		}
//
//		claims, err := jwt.ValidateJWT(refreshToken.Value)
//		if err != nil {
//			log.Error("invalid token: %v", err.Error())
//			http.Error(w, "invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		user, err := getUser.GetUserById(claims.UserId)
//		if err != nil {
//			log.Info("dont find user by id: %v", err.Error())
//			http.Error(w, "invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		accessToken, err := jwt.GenerateAccessToken(user.UserId, user.Role)
//		if err != nil {
//			log.Error("failed to generate access token", err)
//			w.WriteHeader(http.StatusInternalServerError)
//			render.JSON(w, r, resp.Error("failed to generate access token"))
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		render.JSON(w, r, resp.AccessToken(accessToken))
//	}
//}
