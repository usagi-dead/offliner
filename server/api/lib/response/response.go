package response

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
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
