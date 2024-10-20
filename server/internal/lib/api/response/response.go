package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Response godoc
// @Summary Response
// @Description Структура для формата ответа API
// @ID Response
// @Property status string The status of the response (ok or error)
// @Property error string The error message if the status is error (optional)
// @Property access_token string The access token if the status is ok (optional)
type Response struct {
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func AccessToken(token string) Response {
	return Response{
		Status:      StatusOK,
		AccessToken: token,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Email", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
