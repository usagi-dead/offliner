package services

import (
	"golang.org/x/crypto/bcrypt"
	u "server/internal/features/user"
)

func (uus *UserUseCase) SignIn(email string, password string) (string, string, error) {

	user, err := uus.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.HashedPassword), []byte(password)); err != nil {
		return "", "", u.ErrUserNotFound
	}

	if !user.VerifiedEmail {
		return "", "", u.ErrEmailNotConfirmed
	}

	accessToken, err := uus.jwt.GenerateAccessToken(user.UserId, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uus.jwt.GenerateRefreshToken(user.UserId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
