package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
	"net/http"
	"os"
	u "server/internal/features/user"
	"time"
)

var oauthConfigs = map[string]*oauth2.Config{
	"google": &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "http://127.0.0.1:8080/v1/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	},
	"yandex": &oauth2.Config{
		ClientID:     os.Getenv("YANDEX_KEY"),
		ClientSecret: os.Getenv("YANDEX_SECRET"),
		RedirectURL:  "http://127.0.0.1:8080/v1/auth/yandex/callback",
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

func (uuc *UserUseCase) GetAuthURL(provider string) (string, error) {
	config, ok := oauthConfigs[provider]
	if !ok {
		return "", u.ErrUnsupportedProvider
	}

	state := uuid.NewString()
	err := uuc.repo.SaveStateCode(state)
	if err != nil {
		return "", fmt.Errorf("failed to save state: %w", err)
	}

	return config.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

func (uuc *UserUseCase) Callback(provider, state, code string) (bool, string, string, error) {
	config, ok := oauthConfigs[provider]
	if !ok {
		return false, "", "", u.ErrUnsupportedProvider
	}

	isValidState, err := uuc.repo.VerifyStateCode(state)
	if err != nil || !isValidState {
		return false, "", "", u.ErrInvalidState
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return false, "", "", err
	}

	client := config.Client(context.Background(), token)
	user, err := fetchUserInfo(client, provider)
	if err != nil {
		return false, "", "", err
	}

	existingUser, err := uuc.repo.GetUserByEmail(user.Email)
	if errors.Is(err, u.ErrUserNotFound) {
		userID, err := uuc.repo.CreateOauthUser(user)
		if err != nil {
			return false, "", "", err
		}

		accessToken, err := uuc.jwt.GenerateAccessToken(userID, "user")
		if err != nil {
			return false, "", "", err
		}

		refreshToken, err := uuc.jwt.GenerateRefreshToken(userID)
		if err != nil {
			return false, "", "", err
		}

		return false, accessToken, refreshToken, nil
	} else if err != nil {
		return false, "", "", err
	}

	accessToken, err := uuc.jwt.GenerateAccessToken(existingUser.UserId, "user")
	if err != nil {
		return false, "", "", err
	}

	refreshToken, err := uuc.jwt.GenerateRefreshToken(existingUser.UserId)
	if err != nil {
		return false, "", "", err
	}
	return true, accessToken, refreshToken, nil
}

func fetchUserInfo(client *http.Client, provider string) (*u.User, error) {
	var url string
	switch provider {
	case "google":
		url = "https://www.googleapis.com/oauth2/v3/userinfo"
	case "yandex":
		url = "https://login.yandex.ru/info?format=json"
	default:
		return nil, u.ErrUnsupportedProvider
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
		return GoogleToUser(&user), nil
	case "yandex":
		var user YandexUserData
		if err := render.DecodeJSON(resp.Body, &user); err != nil {
			return nil, err
		}
		return YandexToUser(&user), nil
	default:
		return nil, u.ErrUnsupportedProvider
	}
}

func GoogleToUser(googleData *GoogleUserData) *u.User {
	return &u.User{
		Email:         googleData.Email,
		Name:          googleData.Name,
		Surname:       googleData.Surname,
		AvatarUrl:     googleData.AvatarUrl,
		VerifiedEmail: true,
		Role:          "user",
	}
}

func YandexToUser(yandexData *YandexUserData) *u.User {
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

	return &u.User{
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
