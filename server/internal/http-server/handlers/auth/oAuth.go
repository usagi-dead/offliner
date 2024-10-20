package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
	"log/slog"
	"net/http"
	"os"
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

type NotSupportedProviderResponse struct {
	Status string `json:"status" example:"Error"`
	Error  string `json:"error" example:"provider not supported"`
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
// TODO: RELEASE LOGIC
func OauthCallbackHandler(sg StateGenerator, log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
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
			log.Error("Failed to fetch user info: " + err.Error())
			http.Error(w, "Failed to get user info", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, userInfo)
	}
}

func fetchUserInfo(client *http.Client, provider string) (map[string]interface{}, error) {
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

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
