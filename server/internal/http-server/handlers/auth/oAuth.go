package auth

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
	"log/slog"
	"net/http"
	"os"
	"server/internal/lib/api/jwt"
	resp "server/internal/lib/api/response"
	"server/internal/storage/models"
	"time"
)

type StateGenerator interface {
	CreateStateCode() (string, error)
	GetStateCode(stateToken string) (bool, error)
}

var oauthConfigs = map[string]*oauth2.Config{
	"google": &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "http://127.0.0.1:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	},
	"yandex": &oauth2.Config{
		ClientID:     os.Getenv("YANDEX_KEY"),
		ClientSecret: os.Getenv("YANDEX_SECRET"),
		RedirectURL:  "http://127.0.0.1:8080/auth/yandex/callback",
		Endpoint:     yandex.Endpoint,
	},
}

type GoogleUserData struct {
	Email     string  `json:"email"`
	Name      *string `json:"given_name"`
	Surname   *string `json:"family_name"`
	AvatarUrl *string `json:"picture"`
}

type YandexUserData struct {
	Email       string  `json:"default_email"`
	Name        *string `json:"first_name"`
	Surname     *string `json:"last_name"`
	DateOfBirth string  `json:"birthday"`
	Gender      *string `json:"sex"`
	PhoneNumber *string `json:"default_phone.number"`
}

type NotSupportedProviderResponse struct {
	Status string `json:"status" example:"Error"`
	Error  string `json:"error" example:"provider not supported"`
}

type Data interface {
	GetUserByEmail(Email string) (*models.User, error)
	CreateUserAfterOauth(user *models.User) (int64, error)
}

// OauthHandler godoc
// @Summary Start OAuth2 Authorization
// @Tags auth
// @Description Redirects the user to the OAuth provider for authentication.
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider" example("google or yandex")
// @Success 307 "Перенаправление к провайдеру"
// @Failure 404 {object} NotSupportedProviderResponse "Provider not supported"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /auth/{provider} [get]
func OauthHandler(sg StateGenerator, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	log = log.With("op", "OauthHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		config, ok := oauthConfigs[provider]
		if !ok {
			log.Error("not supported provider: " + provider)
			http.Error(w, "provider not supported", http.StatusNotFound)
			return
		}

		state, err := sg.CreateStateCode()
		if err != nil {
			log.Error("error generating state token: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//перенаправление на сторону провайдера
		url := config.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OauthCallbackHandler godoc
func OauthCallbackHandler(d Data, sg StateGenerator, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	log = log.With("op", "OauthCallbackHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		config, ok := oauthConfigs[provider]
		if !ok {
			log.Error("not supported provider: " + provider)
			http.Error(w, "provider not supported", http.StatusNotFound)
			return
		}

		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		ok, err := sg.GetStateCode(state)
		if err != nil {
			log.Error("error getting state: " + err.Error())
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}
		if !ok {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			log.Error("Failed to exchange code: " + err.Error())
			http.Error(w, "Failed to exchange code", http.StatusBadRequest)
			return
		}

		client := config.Client(context.Background(), token)
		userInfo, err := fetchUserInfo(client, provider)
		if err != nil {
			log.Error("failed to fetch user info: " + err.Error())
			http.Error(w, "failed to get user info", http.StatusBadRequest)
			return
		}

		var email string
		switch u := userInfo.(type) {
		case *GoogleUserData:
			email = u.Email
		case *YandexUserData:
			email = u.Email
		default:
			log.Error("unknown user data type")
			http.Error(w, "failed to get user info", http.StatusBadRequest)
			return
		}

		user, err := d.GetUserByEmail(email)
		if err != nil {
			log.Error("database error: " + err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if user != nil {
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

			log.Info("sign in success", slog.String("email", user.Email))

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
			return
		}

		switch u := userInfo.(type) {
		case *GoogleUserData:
			user = ConvertGoogleUserDataToUser(u)
		case *YandexUserData:
			user = ConvertYandexUserDataToUser(u)
		default:
			log.Error("unknown user data type")
			http.Error(w, "unknown user data type", http.StatusBadRequest)
			return
		}

		UserID, err := d.CreateUserAfterOauth(user)
		if err != nil {
			log.Error("failed to create user: " + err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		accessToken, err := jwt.GenerateAccessToken(UserID, "user")
		if err != nil {
			log.Error("failed to generate access token", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to generate access token"))
			return
		}

		refreshToken, err := jwt.GenerateRefreshToken(UserID)
		if err != nil {
			log.Error("failed to generate refresh token", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to generate refresh token"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(15 * 24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, resp.AccessToken(accessToken))
		return
	}
}

func fetchUserInfo(client *http.Client, provider string) (interface{}, error) {
	var url string
	switch provider {
	case "google":
		url = "https://www.googleapis.com/oauth2/v3/userinfo"
	case "yandex":
		url = "https://login.yandex.ru/info?format=json"
	default:
		return nil, fmt.Errorf("unknown provider: %v", provider)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch provider {
	case "google":
		var user GoogleUserData
		if err := render.DecodeJSON(resp.Body, &user); err != nil {
			return nil, err
		}
		return &user, nil
	case "yandex":
		var user YandexUserData
		if err := render.DecodeJSON(resp.Body, &user); err != nil {
			return nil, err
		}
		return &user, nil
	default:
		return nil, fmt.Errorf("unknown provider: %v", provider)
	}
}

func ConvertGoogleUserDataToUser(googleData *GoogleUserData) *models.User {
	return &models.User{
		Email:         googleData.Email,
		Name:          googleData.Name,
		Surname:       googleData.Surname,
		AvatarUrl:     googleData.AvatarUrl,
		VerifiedEmail: true,
		Role:          "user",
	}
}

func ConvertYandexUserDataToUser(yandexData *YandexUserData) *models.User {
	var dob *time.Time
	if yandexData.DateOfBirth != "" {
		parsedDOB, err := time.Parse(time.RFC3339, yandexData.DateOfBirth)
		if err == nil {
			dateOnly := parsedDOB.Format("2006-01-02")
			parsedDateOnly, err := time.Parse("2006-01-02", dateOnly)
			if err == nil {
				dob = &parsedDateOnly
			}
		}
	}

	return &models.User{
		Email:         yandexData.Email,
		Name:          yandexData.Name,
		Surname:       yandexData.Surname,
		DateOfBirth:   dob,
		PhoneNumber:   yandexData.PhoneNumber,
		Gender:        yandexData.Gender,
		VerifiedEmail: true,
		Role:          "user",
	}
}
