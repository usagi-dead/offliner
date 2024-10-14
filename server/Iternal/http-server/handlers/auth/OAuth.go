package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
	"log/slog"
	"net/http"
	"os"
)

var oauthConfigs = map[string]*oauth2.Config{
	"google": &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "http://127.0.0.1:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	},
	"ynadex": &oauth2.Config{
		ClientID:     os.Getenv("YANDEX_KEY"),
		ClientSecret: os.Getenv("YANDEX_SECRET"),
		RedirectURL:  "http://127.0.0.1:8080/auth/ynadex/callback",
		Endpoint:     yandex.Endpoint,
	},
}

var stateStore = map[string]bool{}

func generateStateToken() (stateToken string) {
	stateToken = uuid.NewString()
	stateStore[stateToken] = true
	return
}

func OauthHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	config, ok := oauthConfigs[provider]
	if !ok {
		http.Error(w, "provider not supported", http.StatusNotFound)
		return
	}

	state := generateStateToken()
	//перенаправление на сторону провайдера
	url := config.AuthCodeURL(state, oauth2.AccessTypeOnline) + "&c"
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OauthCallbackHandler(log *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	log = log.With("op", "OauthCallbackHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		config, ok := oauthConfigs[provider]
		if !ok {
			http.Error(w, "provider not supported", http.StatusNotFound)
			return
		}

		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		if !stateStore[state] {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		delete(stateStore, state)

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
	case "ynadex":
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
