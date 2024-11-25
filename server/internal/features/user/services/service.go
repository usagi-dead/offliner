package services

import (
	"log/slog"
	"server/api/lib/emailsender"
	"server/api/lib/jwt"
	"server/internal/features/user"
)

type UserUseCase struct {
	repo user.UserData
	jwt  jwt.JWTService
	log  *slog.Logger
	ess  emailsender.EmailSenderService
}

func NewUserUseCase(d user.UserData, jwt jwt.JWTService, l *slog.Logger, ess emailsender.EmailSenderService) *UserUseCase {
	return &UserUseCase{
		repo: d,
		jwt:  jwt,
		log:  l,
		ess:  ess,
	}
}
