package handlers

import (
	"github.com/go-playground/validator/v10"
	"log/slog"
	"server/internal/features/user"
)

type UserClient struct {
	log      *slog.Logger
	us       user.UserService
	validate *validator.Validate
}

func NewUserClient(log *slog.Logger, aus user.UserService) user.UserHandler {
	return &UserClient{
		log:      log,
		us:       aus,
		validate: validator.New(),
	}
}
