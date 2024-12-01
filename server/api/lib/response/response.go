package response

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	u "server/internal/features/user"
	"strings"
	"time"
)

// Response represents the general structure of an API response
// @Description Structure for a standard API response
type Response struct {
	Status string      `json:"status" example:"success/error"`
	Error  string      `json:"error,omitempty" example:"any error"`
	Data   interface{} `json:"data,omitempty"`
}

const (
	StatusOK    = "success"
	StatusError = "error"
)

type AccessTokenData struct {
	AccessToken string `json:"access_token"`
}

type UserProfileData struct {
	Email       string     `json:"email"`
	Surname     *string    `json:"surname,omitempty"`
	Name        *string    `json:"username,omitempty"`
	Patronymic  *string    `json:"patronymic,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	AvatarUrl   *string    `json:"avatar_url"`
	Gender      *string    `json:"gender,omitempty"`
}

func UserProfile(user *u.User) Response {
	return Response{
		Status: StatusOK,
		Data: UserProfileData{
			Email:       user.Email,
			Surname:     user.Surname,
			Name:        user.Name,
			Patronymic:  user.Patronymic,
			DateOfBirth: user.DateOfBirth,
			PhoneNumber: user.PhoneNumber,
			AvatarUrl:   user.AvatarUrl,
			Gender:      user.Gender,
		},
	}
}

func AccessToken(token string) Response {
	return Response{
		Status: StatusOK,
		Data: AccessTokenData{
			AccessToken: token,
		},
	}
}

func Error(error string) Response {
	return Response{
		Status: StatusError,
		Error:  error,
	}
}

func ValidationError(err error) Response {
	var errMsgs []string

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, err := range validationErrs {
			switch err.ActualTag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
			case "email":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Email", err.Field()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("field %s has an invalid value", err.Field()))
			}
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}
