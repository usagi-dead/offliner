package user

import "errors"

var (
	ErrEmailExists = errors.New("user with this email already exists")
	ErrInternal    = errors.New("internal server error")
)

var (
	ErrInvalidResolutionAvatar = errors.New("invalid resolution avatar(must be 1x1)")
	ErrInvalidTypeAvatar       = errors.New("invalid type avatar")
	ErrSizeAvatar              = errors.New("size of avatar image is large(required <1MB)")
	ErrInvalidState            = errors.New("invalid state")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailNotConfirmed       = errors.New("email not confirmed")
	ErrUnsupportedProvider     = errors.New("unsupported provider")
	ErrEmailAlreadyConfirmed   = errors.New("email already confirmed")
	ErrInvalidConfirmCode      = errors.New("invalid confirm code")
	ErrNoAccessToken           = errors.New("no access token")
	ErrNoRefreshToken          = errors.New("no refresh token")
	ErrExpiredToken            = errors.New("token expired")
	ErrInvalidToken            = errors.New("invalid token")
)
