package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

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
		case "datetime":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid DateOfBirth", err.Field()))
		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Gender", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
