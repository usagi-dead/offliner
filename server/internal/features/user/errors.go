package user

import "errors"

var (
	ErrEmailExists = errors.New("user with this email already exists")
	ErrInternal    = errors.New("internal server error")
)

var (
	ErrInvalidState        = errors.New("invalid state")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailNotConfirmed   = errors.New("email not confirmed")
	ErrUnsupportedProvider = errors.New("unsupported provider")
	ErrDatabase            = errors.New("database error")
)
