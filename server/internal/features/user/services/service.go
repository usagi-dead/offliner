package services

import (
	"log/slog"
	"server/api/lib/jwt"
	"server/internal/features/user"
)

type UserUseCase struct {
	repo user.UserData
	jwt  jwt.JWTService
	log  *slog.Logger
}

func NewUserUseCase(d user.UserData, jwt jwt.JWTService, l *slog.Logger) *UserUseCase {
	return &UserUseCase{
		repo: d,
		jwt:  jwt,
		log:  l,
	}
}
